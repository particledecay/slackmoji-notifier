package slack

import (
	"github.com/slack-go/slack"
)

// Attachment represents an attachment to a Slack message
type Attachment struct {
	ImageURL string
	Text     string
}

// MessageContent represents the content of a Slack message
type MessageContent struct {
	Text        string
	Attachments []Attachment
}

// SendMessage sends a message to the specified Slack channel
func (c *Client) SendMessage(content MessageContent) error {
	_, _, err := c.api.PostMessage(
		c.channel,
		slack.MsgOptionText(content.Text, false),
		slack.MsgOptionAttachments(slack.Attachment{
			ImageURL: content.Attachments[0].ImageURL,
			Text:     content.Attachments[0].Text,
		}),
	)
	return err
}
