package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	Client       *openai.Client
	GPTModel     string
	MaxTokens    int
	SystemPrompt string
}

// NewClient creates a new ChatGPT client
func NewClient(apiKey, gptModel, systemPrompt string, maxTokens int) *Client {
	log.Info().Msgf("creating ChatGPT client with model %s with %d max tokens", gptModel, maxTokens)

	return &Client{
		Client:       openai.NewClient(apiKey),
		GPTModel:     gptModel,
		MaxTokens:    maxTokens,
		SystemPrompt: systemPrompt,
	}
}

// GenerateCompletion sends a list of messages to the ChatGPT API and returns the response
func (c *Client) GenerateCompletion(message string, streamToStdout bool) (string, error) {
	ctx := context.Background()
	var sb strings.Builder

	// Start with our initial prompt
	messagesList := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: c.SystemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:     c.GPTModel,
		MaxTokens: c.MaxTokens,
		Stream:    true,
		Messages:  messagesList,
	}

	stream, err := c.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create chat completion stream")
		return "", err
	}
	defer func(stream *openai.ChatCompletionStream) {
		streamErr := stream.Close()
		if streamErr != nil {
			log.Error().Err(streamErr).Msg("Failed to close created chat completion stream")
		}
	}(stream)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Error().Err(err).Msg("Error while receiving data from stream")
			return "", err
		}

		if streamToStdout {
			fmt.Print(response.Choices[0].Delta.Content)
		} else {
			sb.WriteString(response.Choices[0].Delta.Content)
		}
	}

	if streamToStdout {
		return "", nil
	}

	return sb.String(), nil
}
