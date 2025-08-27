package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

var defaultPrompt = `Generate an edgy, short sentence in modern Gen-Z tone about the given emoji name,
					 and attempt to use a modern and humorous pop culture reference. Do not use proper
					 punctuation, especially periods. Make sure to wrap the exact emoji name as-provided
					 in colons so it can be properly formatted into a Slack emoji. For example, if the
					 emoji name is "smile", the included string should be ":smile:". Don't use other emojis.`

type Config struct {
	Slack struct {
		BotToken string
		AppToken string
		Channel  string
		LogOnly  bool
	}
	OpenAI struct {
		APIKey       string
		Model        string
		MaxTokens    int
		SystemPrompt string
	}
}

func New() *Config {
	config := &Config{}

	log.Debug().Msg("setting Slack configuration")
	config.Slack.BotToken = os.Getenv("SLACK_BOT_TOKEN")
	config.Slack.AppToken = os.Getenv("SLACK_APP_TOKEN")
	config.Slack.Channel = os.Getenv("SLACK_CHANNEL")
	logOnlyValue := os.Getenv("SLACK_LOG_ONLY")
	if logOnlyValue == "" {
		logOnlyValue = "false"
	}
	logOnly, _ := strconv.ParseBool(logOnlyValue)
	config.Slack.LogOnly = logOnly

	log.Debug().Msg("setting OpenAI configuration")
	config.OpenAI.APIKey = os.Getenv("OPENAI_API_KEY")

	config.OpenAI.Model = os.Getenv("OPENAI_MODEL")
	if config.OpenAI.Model == "" {
		log.Info().Msg("OpenAI Model not set, using default")
		config.OpenAI.Model = "gpt-5-nano"
	}

	var maxTokens int
	if maxTokensStr := os.Getenv("OPENAI_MAX_TOKENS"); maxTokensStr != "" {
		parsedMaxTokens, err := strconv.Atoi(maxTokensStr)
		if err != nil {
			log.Warn().Err(err).Msg("error parsing max tokens, using default")
		} else {
			maxTokens = parsedMaxTokens
		}
	}
	config.OpenAI.MaxTokens = maxTokens
	if config.OpenAI.MaxTokens == 0 {
		log.Info().Msg("OpenAI MaxTokens not set, using default")
		config.OpenAI.MaxTokens = 1024
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
