package cmd

import (
	"github.com/particledecay/slackmoji-notifier/pkg/build"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version information`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			if err := build.PrintLongVersion(); err != nil {
				return err
			}
		} else {
			build.PrintVersion()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
