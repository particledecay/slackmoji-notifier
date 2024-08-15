package slack

import (
	"github.com/slack-go/slack"
)

type MessageContent struct {
	Text     string
	ImageURL string
}

func (c *Client) SendMessage(content MessageContent) error {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", content.Text, false, false),
			nil,
			nil,
		),
	}

	if content.ImageURL != "" {
		blocks = append(blocks, slack.NewImageBlock(content.ImageURL, "New Emoji", "", nil))
	}

	if c.useWebhook {
		return slack.PostWebhook(c.webhookURL, &slack.WebhookMessage{
			Blocks: &slack.Blocks{BlockSet: blocks},
		})
	}

	_, _, err := c.api.PostMessage(c.channel, slack.MsgOptionBlocks(blocks...))
	return err
}
