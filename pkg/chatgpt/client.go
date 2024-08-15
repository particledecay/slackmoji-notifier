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
	SystemPrompt string
}

func NewClient(apiKey, gptModel, systemPrompt string) *Client {
	return &Client{
		Client:       openai.NewClient(apiKey),
		GPTModel:     gptModel,
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
		MaxTokens: 1024,
		Stream:    true,
		Messages:  messagesList,
	}

	stream, err := c.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create chat completion stream")
		return "", err
	}
	defer stream.Close()

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
			fmt.Printf(response.Choices[0].Delta.Content)
		} else {
			sb.WriteString(response.Choices[0].Delta.Content)
		}
	}

	if streamToStdout {
		return "", nil
	}

	return sb.String(), nil
}
