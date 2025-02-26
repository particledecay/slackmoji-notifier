package slack

import (
	"errors"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// Client implements the ClientInterface for interacting with the Slack API
type Client struct {
	api          *slack.Client
	socketClient *socketmode.Client
	channel      string
	eventHandler EventHandler
	stopChan     chan struct{}
}

type ClientOption func(*Client)

func NewClient(options ...ClientOption) (ClientInterface, error) {
	client := &Client{}

	for _, option := range options {
		option(client)
	}

	if client.api == nil {
		return nil, errors.New("slack API client must be provided")
	}

	if client.socketClient == nil {
		return nil, errors.New("slack socket mode client must be provided")
	}

	if client.channel == "" {
		return nil, errors.New("channel name must be provided")
	}

	return client, nil
}

func WithAPIToken(botToken, appToken string) ClientOption {
	return func(c *Client) {
		c.api = slack.New(botToken, slack.OptionAppLevelToken(appToken))
		c.socketClient = socketmode.New(
			c.api,
			socketmode.OptionDebug(false),
			socketmode.OptionLog(nil),
		)
	}
}

func WithChannel(channel string) ClientOption {
	return func(c *Client) {
		c.channel = channel
	}
}

func WithEventHandler(handler EventHandler) ClientOption {
	return func(c *Client) {
		c.eventHandler = handler
	}
}
