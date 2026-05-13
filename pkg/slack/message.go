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
	opts := []slack.MsgOption{slack.MsgOptionText(content.Text, false)}
	if len(content.Attachments) > 0 {
		opts = append(opts, slack.MsgOptionAttachments(slack.Attachment{
			ImageURL: content.Attachments[0].ImageURL,
			Text:     content.Attachments[0].Text,
		}))
	}
	_, _, err := c.api.PostMessage(c.channel, opts...)
	return err
}
