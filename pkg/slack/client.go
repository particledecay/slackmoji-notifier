package slack

import (
	"errors"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type Client struct {
	api          *slack.Client
	socketClient *socketmode.Client
	webhookURL   string
	channel      string
	useWebhook   bool
	eventHandler func(event interface{})
}

type ClientOption func(*Client)

func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{}

	for _, option := range options {
		option(client)
	}

	if client.api == nil && client.webhookURL == "" {
		return nil, errors.New("either Slack API token or webhook URL must be provided")
	}

	if client.channel == "" {
		return nil, errors.New("Slack channel must be provided")
	}

	client.useWebhook = client.webhookURL != ""

	return client, nil
}

func WithAPIToken(token string) ClientOption {
	return func(c *Client) {
		c.api = slack.New(token, slack.OptionAppLevelToken(token))
	}
}

func WithAppToken(token string) ClientOption {
	return func(c *Client) {
		if c.api != nil {
			c.socketClient = socketmode.New(
				c.api,
				socketmode.OptionDialer(nil),
			)
		}
	}
}

func WithWebhookURL(url string) ClientOption {
	return func(c *Client) {
		c.webhookURL = url
	}
}

func WithChannel(channel string) ClientOption {
	return func(c *Client) {
		c.channel = channel
	}
}

func WithEventHandler(handler func(event interface{})) ClientOption {
	return func(c *Client) {
		c.eventHandler = handler
	}
}
