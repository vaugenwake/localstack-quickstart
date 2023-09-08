package handlers

import (
	"fmt"
	"localstack-quickstart/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Handler struct {
	Options config.S3Options
	Service *s3.S3
}

func (h *S3Handler) Run() error {
	fmt.Printf("\nCreating S3 bucket: %s\n", h.Options.Name)

	_, err := h.Service.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(h.Options.Name),
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nBucket: '%v' created", h.Options.Name)

	return nil
}
