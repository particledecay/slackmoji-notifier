package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	slackgo "github.com/slack-go/slack"

	"github.com/particledecay/slackmoji-notifier/pkg/chatgpt"
	"github.com/particledecay/slackmoji-notifier/pkg/config"
	"github.com/particledecay/slackmoji-notifier/pkg/slack"
)

func main() {
	cfg := config.New()
	cfg.SetupLogger()

	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid configuration")
	}

	// Initialize ChatGPT client
	chatGPTClient := chatgpt.NewClient(cfg.OpenAI.APIKey, cfg.OpenAI.Model, cfg.OpenAI.SystemPrompt)

	// Initialize Slack client
	var slackClient *slack.Client
	var err error

	eventHandler := func(event interface{}) {
		if ev, ok := event.(*slackgo.EmojiChangedEvent); ok {
			if ev.Type == "add" {
				log.Info().Str("emoji", ev.Name).Msg("New emoji added")

				// Generate sentence with ChatGPT
				prompt := "Generate a fun, short sentence in modern Gen-Z tone using the new emoji :" + ev.Name + ":, making absolutely sure to include the custom emoji as-written in the sentence (and no other emojis)."
				sentence, err := chatGPTClient.GenerateCompletion(prompt, false)
				if err != nil {
					log.Error().Err(err).Msg("Failed to generate sentence")
					return
				}

				// Prepare message content
				messageText := "New emoji added: :" + ev.Name + ":\n\n" + sentence
				messageContent := slack.MessageContent{
					Text:     messageText,
					ImageURL: ev.Value, // This should be the URL of the full-size emoji image
				}

				// Send message to Slack
				if err := slackClient.SendMessage(messageContent); err != nil {
					log.Error().Err(err).Msg("Failed to send message to Slack")
				}
			}
		}
	}

	slackClient, err = slack.NewClient(
		slack.WithAPIToken(cfg.Slack.BotToken),
		slack.WithAppToken(cfg.Slack.AppToken),
		slack.WithChannel(cfg.Slack.Channel),
		slack.WithEventHandler(eventHandler),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Slack client")
	}

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Info().Msg("Received shutdown signal. Shutting down gracefully...")
		cancel()
	}()

	log.Info().Msg("Slackmoji Notifier started")

	// Start listening for Slack events
	go func() {
		if err := slackClient.ListenForEvents(); err != nil {
			log.Error().Err(err).Msg("Slack event listener stopped")
			cancel() // Cancel the context to initiate shutdown
		}
	}()

	// Wait for the context to be cancelled (i.e., for the shutdown signal)
	<-ctx.Done()

	log.Info().Msg("Shutting down Slackmoji Notifier")
}
