package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Slack struct {
		WebhookURL string
		BotToken   string
		AppToken   string
		Channel    string
	}
	OpenAI struct {
		APIKey       string
		Model        string
		SystemPrompt string
	}
	LogLevel zerolog.Level
}

func New() *Config {
	config := &Config{
		LogLevel: zerolog.InfoLevel,
	}

	config.Slack.WebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
	config.Slack.BotToken = os.Getenv("SLACK_BOT_TOKEN")
	config.Slack.AppToken = os.Getenv("SLACK_APP_TOKEN")
	config.Slack.Channel = os.Getenv("SLACK_CHANNEL")

	config.OpenAI.APIKey = os.Getenv("OPENAI_API_KEY")
	config.OpenAI.Model = os.Getenv("OPENAI_MODEL")
	if config.OpenAI.Model == "" {
		config.OpenAI.Model = "gpt-3.5-turbo" // default model
	}
	config.OpenAI.SystemPrompt = os.Getenv("OPENAI_SYSTEM_PROMPT")
	if config.OpenAI.SystemPrompt == "" {
		config.OpenAI.SystemPrompt = "You are a helpful assistant that generates fun, short sentences using given emojis."
	}

	if logLevel, err := zerolog.ParseLevel(strings.ToLower(os.Getenv("LOG_LEVEL"))); err == nil {
		config.LogLevel = logLevel
	}

	return config
}

func (c *Config) SetupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(c.LogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05Z07:00"})
}

func (c *Config) Validate() error {
	if c.Slack.WebhookURL == "" && (c.Slack.BotToken == "" || c.Slack.AppToken == "") {
		return log.Error().Msg("Slack configuration is incomplete. Either SLACK_WEBHOOK_URL or both SLACK_BOT_TOKEN and SLACK_APP_TOKEN must be set")
	}
	if c.Slack.Channel == "" {
		return log.Error().Msg("SLACK_CHANNEL is not set")
	}
	if c.OpenAI.APIKey == "" {
		return log.Error().Msg("OPENAI_API_KEY is not set")
	}
	return nil
}
