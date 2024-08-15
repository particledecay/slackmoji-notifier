package slack

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func (c *Client) ListenForEvents() error {
	if c.useWebhook {
		return errors.New("event listening is not supported when using webhooks")
	}

	if c.socketClient == nil {
		return errors.New("socket mode client is not initialized")
	}

	go func() {
		for evt := range c.socketClient.Events {
			switch evt.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Error().Msg("Could not type cast the event to EventsAPIEvent")
					continue
				}
				c.socketClient.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.EmojiChangedEvent:
						if c.eventHandler != nil {
							c.eventHandler(ev)
						}
					}
				}
			}
		}
	}()

	return c.socketClient.Run()
}
