package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "localstack-config",
	Short: "Localstack config CLI application",
	Long:  "Localstack config CLI helps you to standup and keep your local dev environment setup and work across starts/stops.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "Config file to be used")
}
