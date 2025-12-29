package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

// LLMClient defines the methods that an LLM client should implement
type LLMClient interface {
	GenerateCompletion(ctx context.Context, message string, streamToStdout bool) (string, error)
}

// generateContentWithLLM is a helper function that handles the common logic for generating content
func generateContentWithLLM(ctx context.Context, llm interface{ GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) }, systemPrompt, message string, maxTokens int, streamToStdout bool, providerName string) (string, error) {
	var sb strings.Builder

	messageContents := []llms.MessageContent{
		{Parts: []llms.ContentPart{llms.TextContent{Text: systemPrompt}}, Role: llms.ChatMessageTypeSystem},
		{Parts: []llms.ContentPart{llms.TextContent{Text: message}}, Role: llms.ChatMessageTypeHuman},
	}

	options := []llms.CallOption{
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if streamToStdout {
				fmt.Print(string(chunk))
			} else {
				sb.Write(chunk)
			}
			return nil
		}),
	}

	if maxTokens > 0 {
		options = append(options, llms.WithMaxTokens(maxTokens))
	}

	content, err := llm.GenerateContent(ctx, messageContents, options...)
	if err != nil {
		return "", fmt.Errorf("failed to generate content from %s: %w", providerName, err)
	}

	if streamToStdout {
		return "", nil
	}

	return content.Choices[0].Content, nil
}

// OpenAIClient implements LLMClient for OpenAI models
type OpenAIClient struct {
	llm          *openai.LLM
	modelName    string
	maxTokens    int
	systemPrompt string
}

// NewOpenAIClient creates a new OpenAI LLM client
func NewOpenAIClient(apiKey, modelName, systemPrompt string, maxTokens int) (*OpenAIClient, error) {
	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithModel(modelName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}
	return &OpenAIClient{
		llm:          llm,
		modelName:    modelName,
		maxTokens:    maxTokens,
		systemPrompt: systemPrompt,
	}, nil
}

// GenerateCompletion sends a list of messages to the OpenAI API and returns the response
func (c *OpenAIClient) GenerateCompletion(ctx context.Context, message string, streamToStdout bool) (string, error) {
	return generateContentWithLLM(ctx, c.llm, c.systemPrompt, message, c.maxTokens, streamToStdout, "OpenAI")
}

// OllamaClient implements LLMClient for Ollama models
type OllamaClient struct {
	llm           *ollama.LLM
	modelName     string
	ollamaBaseURL string
	systemPrompt  string
}

// NewOllamaClient creates a new Ollama LLM client
func NewOllamaClient(modelName, ollamaBaseURL, systemPrompt string) (*OllamaClient, error) {
	llm, err := ollama.New(
		ollama.WithModel(modelName),
		ollama.WithServerURL(ollamaBaseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}
	return &OllamaClient{
		llm:           llm,
		modelName:     modelName,
		ollamaBaseURL: ollamaBaseURL,
		systemPrompt:  systemPrompt,
	}, nil
}

// GenerateCompletion sends a list of messages to the Ollama API and returns the response
func (c *OllamaClient) GenerateCompletion(ctx context.Context, message string, streamToStdout bool) (string, error) {
	return generateContentWithLLM(ctx, c.llm, c.systemPrompt, message, 0, streamToStdout, "Ollama")
}

// AnthropicClient implements LLMClient for Anthropic models
type AnthropicClient struct {
	llm          *anthropic.LLM
	modelName    string
	maxTokens    int
	systemPrompt string
}

// NewAnthropicClient creates a new Anthropic LLM client
func NewAnthropicClient(apiKey, modelName, systemPrompt string, maxTokens int) (*AnthropicClient, error) {
	llm, err := anthropic.New(
		anthropic.WithToken(apiKey),
		anthropic.WithModel(modelName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
	}
	return &AnthropicClient{
		llm:          llm,
		modelName:    modelName,
		maxTokens:    maxTokens,
		systemPrompt: systemPrompt,
	}, nil
}

// GenerateCompletion sends a list of messages to the Anthropic API and returns the response
func (c *AnthropicClient) GenerateCompletion(ctx context.Context, message string, streamToStdout bool) (string, error) {
	return generateContentWithLLM(ctx, c.llm, c.systemPrompt, message, c.maxTokens, streamToStdout, "Anthropic")
}

// GoogleAIClient implements LLMClient for Google AI models
type GoogleAIClient struct {
	llm          *googleai.GoogleAI
	modelName    string
	maxTokens    int
	systemPrompt string
}

// NewGoogleAIClient creates a new Google AI LLM client
func NewGoogleAIClient(apiKey, modelName, systemPrompt string, maxTokens int) (*GoogleAIClient, error) {
	llm, err := googleai.New(
		context.Background(),
		googleai.WithAPIKey(apiKey),
		googleai.WithDefaultModel(modelName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GoogleAI client: %w", err)
	}
	return &GoogleAIClient{
		llm:          llm,
		modelName:    modelName,
		maxTokens:    maxTokens,
		systemPrompt: systemPrompt,
	}, nil
}

// GenerateCompletion sends a list of messages to the Google AI API and returns the response
func (c *GoogleAIClient) GenerateCompletion(ctx context.Context, message string, streamToStdout bool) (string, error) {
	return generateContentWithLLM(ctx, c.llm, c.systemPrompt, message, c.maxTokens, streamToStdout, "GoogleAI")
}
