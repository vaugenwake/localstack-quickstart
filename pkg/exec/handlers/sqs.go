package handlers

import (
	"fmt"
	"localstack-quickstart/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSHandler struct {
	Options config.SQSOptions
	Service *sqs.SQS
}

func (h *SQSHandler) Run() error {
	fmt.Printf("\nCreating SQS queue: '%v' \n", h.Options.Name)

	url, err := h.GetQueueUrl()
	if err != nil {
		fmt.Errorf("error checking if queue exists: %v, err: %v", h.Options.Name, err)
	}

	if url != nil {
		fmt.Printf("\nQueue '%v' exists, recreating\n", h.Options.Name)
		h.DeleteQueue(*url)
	}

	attributes := make(map[string]*string)

	if h.Options.VisibilityTimeout != "" {
		attributes["VisibilityTimeout"] = aws.String(h.Options.VisibilityTimeout)
	}

	if h.Options.MessageRetentionPeriod != "" {
		attributes["MessageRetentionPeriod"] = aws.String(h.Options.MessageRetentionPeriod)
	}

	queue, err := h.Service.CreateQueue(&sqs.CreateQueueInput{
		QueueName:  aws.String(h.Options.Name),
		Attributes: attributes,
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nQueue: '%v' created, endpoint: %v\n", h.Options.Name, *queue.QueueUrl)

	return nil
}

func (h *SQSHandler) GetQueueUrl() (*string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: &h.Options.Name,
	}

	queue, err := h.Service.GetQueueUrl(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sqs.ErrCodeQueueDoesNotExist:
				return nil, nil
			}
		}
		return nil, fmt.Errorf("error getting queue url for queue: %v, error: %v", h.Options.Name, err)
	}

	return queue.QueueUrl, nil
}

func (h *SQSHandler) DeleteQueue(queueUrl string) bool {
	input := &sqs.DeleteQueueInput{
		QueueUrl: &queueUrl,
	}

	_, err := h.Service.DeleteQueue(input)
	if err != nil {
		fmt.Errorf("error deleting queue: %v, err: %v", h.Options.Name, err)
		return false
	}

	return true
}
