package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

// LLMClient defines the methods that an LLM client should implement
type LLMClient interface {
	GenerateCompletion(ctx context.Context, message string, streamToStdout bool) (string, error)
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
	var sb strings.Builder

	messageContents := []llms.MessageContent{
		{Parts: []llms.ContentPart{llms.TextContent{Text: c.systemPrompt}}, Role: llms.ChatMessageTypeSystem},
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

	if c.maxTokens > 0 {
		options = append(options, llms.WithMaxTokens(c.maxTokens))
	}

	content, err := c.llm.GenerateContent(ctx, messageContents, options...)
	if err != nil {
		return "", fmt.Errorf("failed to generate content from OpenAI: %w", err)
	}

	if streamToStdout {
		return "", nil
	}

	return content.Choices[0].Content, nil
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
	var sb strings.Builder

	messageContents := []llms.MessageContent{
		{Parts: []llms.ContentPart{llms.TextContent{Text: c.systemPrompt}}, Role: llms.ChatMessageTypeSystem},
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

	content, err := c.llm.GenerateContent(ctx, messageContents, options...)
	if err != nil {
		return "", fmt.Errorf("failed to generate content from Ollama: %w", err)
	}

	if streamToStdout {
		return "", nil
	}

	return content.Choices[0].Content, nil
}
