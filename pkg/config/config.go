package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

const (
	defaultLLMProvider        = "openai"
	defaultOpenAIModel        = "gpt-5-nano"
	defaultOpenAIMaxTokens    = 1024
	defaultOllamaModel        = "llama3.2:1b"
	defaultOllamaBaseURL      = "http://localhost:11434"
	defaultAnthropicModel     = "claude-3.5-haiku"
	defaultAnthropicMaxTokens = 1024
	defaultGoogleAIModel      = "gemini-2.5-flash-lite"
	defaultGoogleAIMaxTokens  = 1024
	defaultSlackLogOnly       = "false"
)

const defaultSystemPrompt = `Generate an edgy, short sentence in modern Gen-Z tone about the given emoji name,
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
		APIKey    string
		Model     string
		MaxTokens int
	}
	Ollama struct {
		Model   string
		BaseURL string
	}
	Anthropic struct {
		APIKey    string
		Model     string
		MaxTokens int
	}
	GoogleAI struct {
		APIKey    string
		Model     string
		MaxTokens int
	}
	LLMProvider  string
	SystemPrompt string
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
	config.LLMProvider = getStringEnvOrDefault("LLM_PROVIDER", defaultLLMProvider)
	config.SystemPrompt = getStringEnvOrDefault("LLM_SYSTEM_PROMPT", defaultSystemPrompt)

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
	case "anthropic":
		setAnthropicConfig(config)
	case "googleai":
		setGoogleAIConfig(config)
	default:
		log.Warn().Str("LLM_PROVIDER", defaultLLMProvider).Msgf("unsupported LLM_PROVIDER: %s, using default", config.LLMProvider)
		config.LLMProvider = defaultLLMProvider
		setOpenAIConfig(config)
	}

	return config
}

func setOpenAIConfig(config *Config) {
	config.OpenAI.APIKey = os.Getenv("OPENAI_API_KEY")
	config.OpenAI.Model = getStringEnvOrDefault("OPENAI_MODEL", defaultOpenAIModel)
	config.OpenAI.MaxTokens = getIntEnvOrDefault("OPENAI_MAX_TOKENS", defaultOpenAIMaxTokens)
	log.Info().Str("model", config.OpenAI.Model).Msg("using OpenAI model")
}

func getStringEnvOrDefault(envVar, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		log.Info().Str(envVar, defaultValue).Msg("environment variable not set, using default")
		return defaultValue
	}
	return value
}

func getIntEnvOrDefault(envVar string, defaultValue int) int {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		log.Info().Int(envVar, defaultValue).Msg("environment variable not set, using default")
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Warn().Err(err).Str(envVar, valueStr).Int("default", defaultValue).Msg("error parsing environment variable, using default")
		return defaultValue
	}

	return parsedValue
}

func setAnthropicConfig(config *Config) {
	config.Anthropic.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	config.Anthropic.Model = getStringEnvOrDefault("ANTHROPIC_MODEL", defaultAnthropicModel)
	config.Anthropic.MaxTokens = getIntEnvOrDefault("ANTHROPIC_MAX_TOKENS", defaultAnthropicMaxTokens)
	log.Info().Str("model", config.Anthropic.Model).Msg("using Anthropic model")
}

func setGoogleAIConfig(config *Config) {
	config.GoogleAI.APIKey = os.Getenv("GOOGLEAI_API_KEY")
	config.GoogleAI.Model = getStringEnvOrDefault("GOOGLEAI_MODEL", defaultGoogleAIModel)
	config.GoogleAI.MaxTokens = getIntEnvOrDefault("GOOGLEAI_MAX_TOKENS", defaultGoogleAIMaxTokens)
	log.Info().Str("model", config.GoogleAI.Model).Msg("using GoogleAI model")
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
	case "anthropic":
		if c.Anthropic.APIKey == "" {
			log.Error().Msg("ANTHROPIC_API_KEY is not set")
			return errors.New("ANTHROPIC_API_KEY is not set")
		}
	case "googleai":
		if c.GoogleAI.APIKey == "" {
			log.Error().Msg("GOOGLEAI_API_KEY is not set")
			return errors.New("GOOGLEAI_API_KEY is not set")
		}
	default:
		return fmt.Errorf("unsupported LLM_PROVIDER: %s", c.LLMProvider)
	}

	return nil
}
