package notifier

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/particledecay/slackmoji-notifier/pkg/chatgpt"
	"github.com/particledecay/slackmoji-notifier/pkg/slack"
)

const eventThreshold = 1 * time.Minute

type Notifier struct {
	slackClient     slack.ClientInterface
	chatGPTClient   chatgpt.ClientInterface
	processedEvents map[string]time.Time
	knownEmojis     map[string]bool
	eventsMutex     sync.Mutex
}

func New(chatGPTClient chatgpt.ClientInterface) *Notifier {
	n := &Notifier{
		chatGPTClient:   chatGPTClient,
		processedEvents: make(map[string]time.Time),
		knownEmojis:     make(map[string]bool),
	}
	n.startCleanupRoutine()
	return n
}

func (n *Notifier) SetSlackClient(client slack.ClientInterface) {
	n.slackClient = client
}

func (n *Notifier) HandleEvent(event interface{}) {
	log.Debug().Interface("event", event).Msgf("notifier received event of type %T", event)

	switch evt := event.(type) {
	case socketmode.Event:
		log.Debug().Msg("event is a socketmode event")
		n.handleSocketModeEvent(evt)
	default:
		log.Debug().Msgf("unhandled event type %T", event)
	}
}

func (n *Notifier) handleSocketModeEvent(event socketmode.Event) {
	log.Debug().Str("type", string(event.Type)).Msg("handling socketmode event")

	if event.Type == socketmode.EventTypeEventsAPI {
		var payload struct {
			EventID      string `json:"event_id"`
			EventTime    int64  `json:"event_time"`
			RetryAttempt int    `json:"retry_attempt"`
		}
		if err := json.Unmarshal(event.Request.Payload, &payload); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal payload")
			return
		}

		eventTime := time.Unix(payload.EventTime, 0)
		if time.Since(eventTime) > eventThreshold {
			log.Debug().Str("event_id", payload.EventID).Time("event_time", eventTime).Msg("ignoring old event")
			return
		}

		if payload.RetryAttempt > 0 {
			log.Debug().Str("event_id", payload.EventID).Int("retry_attempt", payload.RetryAttempt).Msg("ignoring retry event")
			return
		}

		eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
		if !ok {
			log.Debug().Msg("event data is not an EventsAPIEvent")
			return
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.EmojiChangedEvent:
				if ev.Subtype == "add" {
					n.handleNewEmoji(ev.Name, ev.Value)
				}
			default:
				log.Debug().Str("type", innerEvent.Type).Msg("unhandled inner event type")
			}
		}
	} else {
		log.Debug().Str("type", string(event.Type)).Msg("event is not an EventsAPI event, skipping")
	}
}

func (n *Notifier) handleNewEmoji(name, value string) {
	n.eventsMutex.Lock()
	if n.knownEmojis[name] {
		n.eventsMutex.Unlock()
		log.Debug().Str("emoji", name).Msg("ignoring known emoji")
		return
	}
	n.knownEmojis[name] = true
	n.eventsMutex.Unlock()

	log.Info().Str("emoji", name).Msg("handling new emoji")

	sentence, err := n.chatGPTClient.GenerateCompletion("emoji name: "+name, false)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate sentence")
		return
	}

	log.Debug().Str("sentence", sentence).Msg("generated sentence for new emoji")

	messageText := "New emoji: :" + name + ":\n\n" + sentence
	messageContent := slack.MessageContent{
		Text:     messageText,
		ImageURL: value,
	}

	log.Debug().Interface("messageContent", messageContent).Msg("sending message to Slack")

	if err := n.slackClient.SendMessage(messageContent); err != nil {
		log.Error().Err(err).Msg("failed to send message to Slack")
	} else {
		log.Debug().Msg("message sent successfully to Slack")
	}
}

func (n *Notifier) cleanupProcessedEvents() {
	n.eventsMutex.Lock()
	defer n.eventsMutex.Unlock()

	threshold := time.Now().Add(-eventThreshold)
	for id, timestamp := range n.processedEvents {
		if timestamp.Before(threshold) {
			delete(n.processedEvents, id)
		}
	}
	log.Debug().Msg("cleaned up processed events")
}

func (n *Notifier) startCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(eventThreshold)
		defer ticker.Stop()
		for range ticker.C {
			n.cleanupProcessedEvents()
		}
	}()
}
