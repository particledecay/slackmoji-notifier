package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

const (
	defaultLLMProvider     = "openai"
	defaultOpenAIModel     = "gpt-5-nano"
	defaultOpenAIMaxTokens = 1024
	defaultOllamaModel     = "llama3.2:1b"
	defaultOllamaBaseURL   = "http://localhost:11434"
	defaultSlackLogOnly    = "false"
)

const defaultOpenAISystemPrompt = `Generate an edgy, short sentence in modern Gen-Z tone about the given emoji name,
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
	Ollama struct {
		Model   string
		BaseURL string
	}
	LLMProvider string
}

func New() *Config {
	config := &Config{}

	log.Debug().Msg("setting Slack configuration")
	config.Slack.BotToken = os.Getenv("SLACK_BOT_TOKEN")
	config.Slack.AppToken = os.Getenv("SLACK_APP_TOKEN")
	config.Slack.Channel = os.Getenv("SLACK_CHANNEL")
	logOnlyValue := os.Getenv("SLACK_LOG_ONLY")
	if logOnlyValue == "" {
		logOnlyValue = defaultSlackLogOnly
	} else {
		log.Info().Str("SLACK_LOG_ONLY", logOnlyValue).Msg("SLACK_LOG_ONLY explicitly set")
	}
	logOnly, _ := strconv.ParseBool(logOnlyValue)
	config.Slack.LogOnly = logOnly

	log.Debug().Msg("setting LLM configuration")
	config.LLMProvider = os.Getenv("LLM_PROVIDER")
	if config.LLMProvider == "" {
		log.Info().Str("LLM_PROVIDER", defaultLLMProvider).Msg("LLM_PROVIDER not set, using default")
		config.LLMProvider = defaultLLMProvider
	}

	switch config.LLMProvider {
	case "openai":
		setOpenAIConfig(config)
	case "ollama":
		config.Ollama.Model = os.Getenv("OLLAMA_MODEL")
		if config.Ollama.Model == "" {
			log.Info().Str("OLLAMA_MODEL", defaultOllamaModel).Msg("Ollama Model not set, using default")
			config.Ollama.Model = defaultOllamaModel
		}
		config.Ollama.BaseURL = os.Getenv("OLLAMA_BASE_URL")
		if config.Ollama.BaseURL == "" {
			log.Info().Str("OLLAMA_BASE_URL", defaultOllamaBaseURL).Msg("Ollama BaseURL not set, using default")
			config.Ollama.BaseURL = defaultOllamaBaseURL
		}
		log.Info().Str("model", config.Ollama.Model).Str("baseURL", config.Ollama.BaseURL).Msg("using Ollama model")
	default:
		log.Warn().Str("LLM_PROVIDER", defaultLLMProvider).Msgf("unsupported LLM_PROVIDER: %s, using default", config.LLMProvider)
		config.LLMProvider = defaultLLMProvider
		setOpenAIConfig(config)
	}

	return config
}

func setOpenAIConfig(config *Config) {
	config.OpenAI.APIKey = os.Getenv("OPENAI_API_KEY")
	config.OpenAI.Model = os.Getenv("OPENAI_MODEL")
	if config.OpenAI.Model == "" {
		log.Info().Str("OPENAI_MODEL", defaultOpenAIModel).Msg("OpenAI Model not set, using default")
		config.OpenAI.Model = defaultOpenAIModel
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
		log.Info().Int("OPENAI_MAX_TOKENS", defaultOpenAIMaxTokens).Msg("OpenAI MaxTokens not set, using default")
		config.OpenAI.MaxTokens = defaultOpenAIMaxTokens
	}

	log.Info().Str("model", config.OpenAI.Model).Msg("using OpenAI model")
	config.OpenAI.SystemPrompt = os.Getenv("OPENAI_SYSTEM_PROMPT")
	if config.OpenAI.SystemPrompt == "" {
		log.Debug().Str("OPENAI_SYSTEM_PROMPT", defaultOpenAISystemPrompt).Msg("OpenAI System Prompt not set, using default")
		config.OpenAI.SystemPrompt = defaultOpenAISystemPrompt
	}
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

	switch c.LLMProvider {
	case "openai":
		if c.OpenAI.APIKey == "" {
			log.Error().Msg("OPENAI_API_KEY is not set")
			return errors.New("OPENAI_API_KEY is not set")
		}
	case "ollama":
		if c.Ollama.Model == "" {
			log.Error().Msg("OLLAMA_MODEL is not set")
			return errors.New("OLLAMA_MODEL is not set")
		}
		if c.Ollama.BaseURL == "" {
			log.Error().Msg("OLLAMA_BASE_URL is not set")
			return errors.New("OLLAMA_BASE_URL is not set")
		}
	default:
		return fmt.Errorf("unsupported LLM_PROVIDER: %s", c.LLMProvider)
	}

	return nil
}
