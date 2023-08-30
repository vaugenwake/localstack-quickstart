package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"localstack-quickstart/config"
)

type SQSHandler struct {
	Options config.SQSOptions
	session *session.Session
}

func (h *SQSHandler) SetSession(s *session.Session) {
	h.session = s
}

func (h *SQSHandler) Run() error {
	fmt.Printf("\nCreating SQS queue: '%v'", h.Options.Name)

	srv := sqs.New(h.session)

	queue, err := srv.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(h.Options.Name),
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nQueue: '%v' created, endpoint: %v", h.Options.Name, *queue.QueueUrl)

	return nil
}
