package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "slackmoji-notifier",
		Short: "A CLI tool to notify users of new emojis in a Slack workspace",
		Long:  `A CLI tool to notify users of new emojis in a Slack workspace.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupLogger()
		},
	}

	verbose bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute root command")
		os.Exit(1)
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Debug().Msg("verbose mode enabled")
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
}
