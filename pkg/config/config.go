package config

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

var defaultPrompt = `Generate an edgy, short sentence in modern Gen-Z tone about the given emoji name.
					 Surround the name with colons (e.g., if 'smiling' is the name, then :smiling: should
					 be in the sentence), making absolutely sure to include the custom emoji as-written in
					 the sentence (and you must not include any other emojis). The custom emoji must
					 be the only emoji in the sentence and cannot be specified without wrapping colons.
					 You cannot generate a sentence without an emoji. Think more "edgy" and less "corny".
					 The sentence should attempt to be about the emoji, using the emoji name to make
					 assumptions about its meaning. Do not remove any characters from the emoji name,
					 otherwise it will not display properly in Slack.`

type Config struct {
	Slack struct {
		BotToken string
		AppToken string
		Channel  string
	}
	OpenAI struct {
		APIKey       string
		Model        string
		SystemPrompt string
	}
}

func New() *Config {
	config := &Config{}

	log.Debug().Msg("setting Slack configuration")
	config.Slack.BotToken = os.Getenv("SLACK_BOT_TOKEN")
	config.Slack.AppToken = os.Getenv("SLACK_APP_TOKEN")
	config.Slack.Channel = os.Getenv("SLACK_CHANNEL")

	log.Debug().Msg("setting OpenAI configuration")
	config.OpenAI.APIKey = os.Getenv("OPENAI_API_KEY")
	config.OpenAI.Model = os.Getenv("OPENAI_MODEL")
	if config.OpenAI.Model == "" {
		log.Debug().Msg("OpenAI Model not set, using default")
		config.OpenAI.Model = "gpt-3.5-turbo"
	}
	log.Debug().Str("model", config.OpenAI.Model).Msg("using OpenAI model")
	config.OpenAI.SystemPrompt = os.Getenv("OPENAI_SYSTEM_PROMPT")
	if config.OpenAI.SystemPrompt == "" {
		log.Debug().Msg("OpenAI System Prompt not set, using default")
		config.OpenAI.SystemPrompt = defaultPrompt
	}

	return config
}

func (c *Config) Validate() error {
	if c.Slack.BotToken == "" {
		log.Error().Msg("SLACK_BOT_TOKEN is not set")
		return errors.New("SLACK_BOT_TOKEN is not set")
	}
	if c.Slack.AppToken == "" {
		log.Error().Msg("SLACK_APP_TOKEN is not set")
		return errors.New("SLACK_APP_TOKEN is not set")
	}
	if c.Slack.Channel == "" {
		log.Error().Msg("SLACK_CHANNEL is not set")
		return errors.New("SLACK_CHANNEL is not set")
	}
	if c.OpenAI.APIKey == "" {
		log.Error().Msg("OPENAI_API_KEY is not set")
		return errors.New("OPENAI_API_KEY is not set")
	}
	return nil
}
