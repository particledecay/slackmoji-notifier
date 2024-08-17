package slack

import (
	"github.com/slack-go/slack"
)

// MessageContent represents the content of a Slack message
type MessageContent struct {
	Text     string
	ImageURL string
}

// SendMessage sends a message to the specified Slack channel
func (c *Client) SendMessage(content MessageContent) error {
	_, _, err := c.api.PostMessage(
		c.channel,
		slack.MsgOptionText(content.Text, false),
		slack.MsgOptionAttachments(slack.Attachment{
			ImageURL: content.ImageURL,
		}),
	)
	return err
}
