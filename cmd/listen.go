package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/slackmoji-notifier/internal/notifier"
	"github.com/particledecay/slackmoji-notifier/pkg/chatgpt"
	"github.com/particledecay/slackmoji-notifier/pkg/config"
	"github.com/particledecay/slackmoji-notifier/pkg/slack"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for new emoji events in Slack",
	Long:  `Start listening for new emoji events in the configured Slack workspace and send notifications.`,
	Run:   runListen,
}

func init() {
	rootCmd.AddCommand(listenCmd)
}

func runListen(cmd *cobra.Command, args []string) {
	log.Debug().Msg("starting listen command")

	cfg := config.New()
	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("invalid configuration")
	}
	log.Debug().Msg("configuration validated successfully")

	chatGPTClient := chatgpt.NewClient(cfg.OpenAI.APIKey, cfg.OpenAI.Model, cfg.OpenAI.SystemPrompt, cfg.OpenAI.MaxTokens)
	log.Debug().Msg("ChatGPT client initialized")

	n := notifier.New(chatGPTClient, cfg.Slack.LogOnly)
	log.Debug().Msg("notifier created")

	debugEventHandler := func(event interface{}) {
		log.Debug().Interface("event", event).Msg("received Slack event")
		n.HandleEvent(event)
	}

	log.Debug().Str("channel", cfg.Slack.Channel).Msg("initializing Slack client")
	slackClient, err := slack.NewClient(
		slack.WithAPIToken(cfg.Slack.BotToken, cfg.Slack.AppToken),
		slack.WithChannel(cfg.Slack.Channel),
		slack.WithEventHandler(debugEventHandler),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Slack client")
	}
	log.Debug().Msg("Slack client created successfully")

	n.SetSlackClient(slackClient)
	log.Debug().Msg("Slack client set in notifier")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("received shutdown signal, shutting down gracefully...")
		cancel()
	}()

	log.Info().Msg("starting event listener")

	if err := slackClient.ListenForEvents(); err != nil {
		log.Error().Err(err).Msg("event listener stopped")
		cancel()
	}

	<-ctx.Done()
	log.Info().Msg("shutting down")

	// Give ongoing operations some time to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	slackClient.Stop()
	log.Debug().Msg("Slack client stopped")

	select {
	case <-shutdownCtx.Done():
		log.Warn().Msg("shutdown timed out")
	default:
		log.Info().Msg("shutdown completed successfully")
	}
}
