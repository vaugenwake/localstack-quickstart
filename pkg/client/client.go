package client

import (
	"context"
	"fmt"
	"localstack-quickstart/pkg/config"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	Connection *config.Connection
}

func (c *Client) Connect() (*session.Session, error) {
	awsConfig := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String(c.Connection.GetEndpoint()),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("error setting up session: %v", err)
	}

	return sess, nil
}

func (c *Client) HealthCheck(sess *session.Session) bool {
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
