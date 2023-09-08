package cmd

import (
	"context"
	"fmt"
	"localstack-quickstart/pkg/client"
	"localstack-quickstart/pkg/config"
	"localstack-quickstart/pkg/errors"
	"localstack-quickstart/pkg/exec"
	"time"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize and setup application with config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("--config= parameter not provided")
		}

		parsedConfig, err := config.ParseConfigFile(configFile)
		if err != nil {
			return fmt.Errorf("could not parse config file: %v, error: %v", configFile, err.Error())
		}

		client := client.Client{
			Connection: &parsedConfig.Connection,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		sess, err := client.InitSession()
		if err != nil {
			return fmt.Errorf("error estabilishing session: %v", err.Error())
		}

		if !client.HealthCheck(&ctx, sess) {
			return fmt.Errorf("could not connect, retry limit reached")
		}

		errorCollector := errors.ErrorsBag{}

		executor := &exec.ExecutionPlan{}

		err = executor.Plan(&parsedConfig.Resources, sess)
		if err != nil {
			return fmt.Errorf("could not create execution plan, %v", err.Error())
		}

		err = executor.Exec(&errorCollector)
		if err != nil {
			return fmt.Errorf("error running execution plan, %v", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
