package slack

import (
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/socketmode"
)

// EventHandler is a function type for handling Slack events
type EventHandler func(event interface{})

// ListenForEvents starts listening for Slack events
func (c *Client) ListenForEvents() error {
	c.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-c.stopChan:
				return
			default:
				if err := c.socketClient.Run(); err != nil {
					log.Error().Err(err).Msg("failed to run socket client")
					return
				}
			}
		}
	}()

	go func() {
		for evt := range c.socketClient.Events {
			c.handleEvent(evt)
		}
	}()

	return nil
}

// handleEvent processes incoming Slack events
func (c *Client) handleEvent(evt socketmode.Event) {
	if c.eventHandler != nil {
		c.eventHandler(evt)
	}
}

// Stop signals the event listener to stop
func (c *Client) Stop() {
	if c.stopChan != nil {
		close(c.stopChan)
	}
}
