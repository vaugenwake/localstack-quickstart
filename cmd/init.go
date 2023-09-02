package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"localstack-quickstart/errors"
	"localstack-quickstart/pkg/config"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize and setup application with config",
	Run: func(cmd *cobra.Command, args []string) {
		errorCollecter := &errors.ErrorsBag{}

		var configFile string
		if cmd.Flag("config").Value.String() != "" {
			configFile = cmd.Flag("config").Value.String()
		} else {
			configFile = "config.yml"
		}

		parsedConfig, err := config.ParseConfigFile(configFile)
		if err != nil {
			errorCollecter.Add("Fatal", err.Error())
			//printError(errorCollecter)
			os.Exit(1)
		}

		fmt.Println(parsedConfig)

		//sess, err := connectToAws(parsedConfig)
		//if err != nil {
		//	errorCollecter.Add("Fatal", err.Error())
		//	printError(errorCollecter)
		//	os.Exit(1)
		//}
		//
		//if !checkHealthy(sess) {
		//	errorCollecter.Add("Fatal", "Could not connect to localstack, retry limit reached")
		//	printError(errorCollecter)
		//	os.Exit(1)
		//}
		//
		//executor := &exec.ExecutionPlan{}
		//
		//err = executor.Plan(&parsedConfig.Resources, sess)
		//if err != nil {
		//	errorCollecter.Add("Fatal", err.Error())
		//}
		//
		//err = executor.Exec()
		//if err != nil {
		//	errorCollecter.Add("Fatal", err.Error())
		//}
		//
		//printError(errorCollecter)
		//os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
