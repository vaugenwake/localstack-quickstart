package main

import (
	"context"
	"fmt"
	"localstack-quickstart/exec"
	"os"
	"time"

	"localstack-quickstart/config"
	"localstack-quickstart/errors"
	"localstack-quickstart/inputs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jedib0t/go-pretty/table"
)

func connectToAws(config *config.Config) (*session.Session, error) {
	awsConfig := &aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String(config.GetEndpoint()),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Printf("Error setting up connection context: %v", err)
		return nil, err
	}

	return sess, nil
}

func checkHealthy(sess *session.Session) bool {
	retries := 3
	retryInterval := 5 * time.Second

	s3Client := s3.New(sess)

	for i := 0; i < retries; i++ {
		_, err := s3Client.ListBucketsWithContext(context.Background(), &s3.ListBucketsInput{})
		if err == nil {
			return true
		}

		fmt.Printf("Connection Attempt: %d, Session is not healthy: %v\n", i+1, err)

		if i < retries-1 {
			fmt.Printf("Retrying in %v...\n", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	return false
}

func main() {

	errorCollecter := &errors.ErrorsBag{}

	inputs, err := inputs.ParseInputFlags()
	if err != nil {
		panic("Inputs failer")
	}

	parsedConfig, err := config.ParseConfigFile(inputs.ConfigFile)
	if err != nil {
		errorCollecter.Add("Fatal", err.Error())
	}

	executor := &exec.ExecutionPlan{}
	err = executor.Plan(&parsedConfig.Resources)
	if err != nil {
		errorCollecter.Add("Fatal", err.Error())
	}

	err = executor.Exec()
	if err != nil {
		errorCollecter.Add("Fatal", err.Error())
	}

	if errorCollecter.Any() {
		t := table.NewWriter()
		t.SetTitle("Execution Errors")

		t.AppendHeader(table.Row{"#", "Level", "Error"})

		for idx, err := range errorCollecter.All() {
			t.AppendRow(table.Row{idx, err.Level, err.Message})
		}

		fmt.Println(t.Render())
		os.Exit(1)
	}

	// sess, err := connectToAws(config)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// if !checkHealthy(sess) {
	// 	fmt.Println("Could not establish healthy connection to localstack service")
	// 	os.Exit(1)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	// dynamoSrv := dynamodb.New(sess)

	// fmt.Println("Tables:")

	// for {
	// 	result, err := dynamoSrv.ListTablesWithContext(ctx, &dynamodb.ListTablesInput{})
	// 	if err != nil {
	// 		if aerr, ok := err.(awserr.Error); ok {
	// 			switch aerr.Code() {
	// 			case dynamodb.ErrCodeInternalServerError:
	// 				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
	// 			default:
	// 				fmt.Println(aerr.Error())
	// 			}
	// 		} else {
	// 			// Print the error, cast err to awserr.Error to get the Code and
	// 			// Message from an error.
	// 			fmt.Println(err.Error())
	// 		}
	// 		return
	// 	}

	// 	for _, n := range result.TableNames {
	// 		fmt.Println(*n)
	// 	}

	// 	if result.LastEvaluatedTableName == nil {
	// 		break
	// 	}
	// }
}
