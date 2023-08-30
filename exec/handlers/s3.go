package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"localstack-quickstart/config"
)

type S3Handler struct {
	Options config.S3Options
	session *session.Session
}

func (h *S3Handler) SetSession(s *session.Session) {
	h.session = s
}

func (h *S3Handler) Run() error {
	fmt.Printf("\nCreating S3 bucket: %s\n", h.Options.Name)

	srv := s3.New(h.session)

	_, err := srv.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(h.Options.Name),
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nBucket: '%v' created", h.Options.Name)

	return nil
}
