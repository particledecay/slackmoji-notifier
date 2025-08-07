package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/particledecay/slackmoji-notifier/pkg/llm"
	"github.com/particledecay/slackmoji-notifier/pkg/slack"
)

const eventThreshold = 1 * time.Minute

type Notifier struct {
	slackClient     slack.ClientInterface
	llmClient       llm.LLMClient
	processedEvents map[string]time.Time
	knownEmojis     map[string]bool
	eventsMutex     sync.Mutex
	logOnly         bool
}

func New(llmClient llm.LLMClient, logOnly bool) *Notifier {
	n := &Notifier{
		llmClient:       llmClient,
		processedEvents: make(map[string]time.Time),
		knownEmojis:     make(map[string]bool),
		logOnly:         logOnly,
	}
	n.startCleanupRoutine()
	return n
}

func (n *Notifier) SetSlackClient(client slack.ClientInterface) {
	n.slackClient = client
}

func (n *Notifier) HandleEvent(event interface{}) {
	log.Debug().Interface("event", event).Msgf("notifier received event of type %T", event)

	ctx := context.Background()

	switch evt := event.(type) {
	case socketmode.Event:
		log.Debug().Msg("event is a socketmode event")
		n.handleSocketModeEvent(ctx, evt)
	default:
		log.Debug().Msgf("unhandled event type %T", event)
	}
}

type SocketModePayload struct {
	EventID      string `json:"event_id"`
	EventTime    int64  `json:"event_time"`
	RetryAttempt int    `json:"retry_attempt"`
}

func (n *Notifier) handleSocketModeEvent(ctx context.Context, event socketmode.Event) {
	log.Debug().Str("type", string(event.Type)).Msg("handling socketmode event")

	if event.Type == socketmode.EventTypeEventsAPI {
		var payload SocketModePayload
		if err := json.Unmarshal(event.Request.Payload, &payload); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal payload")
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
				// slack sends emoji changes multiple times for some reason so we ignore older ones
				eventTime := time.Unix(payload.EventTime, 0)
				if time.Since(eventTime) > eventThreshold {
					log.Debug().
						Str("event_id", payload.EventID).
						Time("event_time", eventTime).
						Str("emoji", ev.Name).
						Msg("ignoring old event")
					return
				}

				// even within the threshold sometimes we get retries, so we ignore those too
				if payload.RetryAttempt > 0 {
					log.Debug().
						Str("event_id", payload.EventID).
						Int("retry_attempt", payload.RetryAttempt).
						Str("emoji", ev.Name).
						Msg("ignoring retry event")
					return
				}

				// we only care about new emojis
				if ev.Subtype == "add" {
					n.handleNewEmoji(ctx, ev.Name, ev.Value)
				} else if ev.Subtype == "remove" {
					n.handleRemovedEmoji(ev.Name)
				}
			default:
				log.Debug().Str("type", innerEvent.Type).Msg("unhandled inner event type")
			}
		}
	} else {
		log.Debug().Str("type", string(event.Type)).Msg("event is not an EventsAPI event, skipping")
	}
}

func (n *Notifier) handleRemovedEmoji(name string) {
	n.eventsMutex.Lock()
	defer n.eventsMutex.Unlock()

	val, ok := n.knownEmojis[name]
	if !ok {
		log.Debug().Str("emoji", name).Msg("ignoring unaccounted emoji")
		return
	}
	if !val {
		log.Debug().Str("emoji", name).Msg("ignoring already deleted emoji")
		return
	}

	log.Info().Str("emoji", name).Msg("removing emoji from known emojis")
	n.knownEmojis[name] = false
}

func (n *Notifier) handleNewEmoji(ctx context.Context, name, value string) {
	n.eventsMutex.Lock()
	defer n.eventsMutex.Unlock()

	if n.knownEmojis[name] {
		log.Debug().Str("emoji", name).Msg("ignoring known emoji")
		return
	}
	n.knownEmojis[name] = true

	log.Info().Str("emoji", name).Msg("handling new emoji")

	if n.logOnly {
		log.Info().
			Str("emoji", name).
			Msg("logOnly mode: would have generated sentence and sent Slack message")
		return
	}

	sentence, err := n.llmClient.GenerateCompletion(ctx, "emoji name: "+name, false)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate sentence")
		n.knownEmojis[name] = false
		return
	}

	log.Debug().Str("sentence", sentence).Msg("generated sentence for new emoji")

	// Construct the full-size image URL
	fullSizeImageURL := value
	if !strings.Contains(value, "?") {
		fullSizeImageURL += "?size=512"
	} else {
		fullSizeImageURL += "&size=512"
	}

	messageText := fmt.Sprintf("*NEW EMOJI ADDED!*\n*Example Usage:*\n%s", sentence)
	messageContent := slack.MessageContent{
		Text: messageText,
		Attachments: []slack.Attachment{
			{
				ImageURL: fullSizeImageURL,
				Text:     name,
			},
		},
	}

	log.Debug().Interface("messageContent", messageContent).Msg("sending message to Slack")

	if err := n.slackClient.SendMessage(messageContent); err != nil {
		log.Error().Err(err).Msg("failed to send message to Slack")
		n.knownEmojis[name] = false
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
