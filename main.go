package main

import (
	"context"
	"fmt"
	"localstack-quickstart/cmd"
	"localstack-quickstart/pkg/config"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jedib0t/go-pretty/table"
	"localstack-quickstart/errors"
)

func connectToAws(config *config.Config) (*session.Session, error) {
	awsConfig := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(config.GetEndpoint()),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
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

func printError(e *errors.ErrorsBag) {
	if e.Any() {
		t := table.NewWriter()
		t.SetTitle("Execution Errors")

		t.AppendHeader(table.Row{"#", "Level", "Error"})

		for idx, err := range e.All() {
			t.AppendRow(table.Row{idx, err.Level, err.Message})
		}

		fmt.Println(t.Render())
		os.Exit(1)
	}
}

func main() {
	cmd.Execute()
}
