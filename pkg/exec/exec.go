package exec

import (
	"context"
	"fmt"
	"localstack-quickstart/pkg/config"
	"localstack-quickstart/pkg/errors"
	"localstack-quickstart/pkg/exec/handlers"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Handler interface {
	// SetSession(s *session.Session)
	Run() error
}

type ExecutionStep struct {
	Type    config.ResourceType
	handler Handler
}

type ExecutionPlan struct {
	Steps []ExecutionStep
	ctx   *context.Context
}

func handlerFactory(resourceType config.ResourceType, options interface{}, sess *session.Session) (Handler, error) {
	var handler Handler

	switch resourceType {
	case config.S3:
		if s3Opts, ok := options.(config.S3Options); ok {
			handler = &handlers.S3Handler{
				Options: s3Opts,
				Service: s3.New(sess),
			}
		}
	case config.SQS:
		if sqsOpts, ok := options.(config.SQSOptions); ok {
			handler = &handlers.SQSHandler{
				Options: sqsOpts,
				Service: sqs.New(sess),
			}
		}
	}

	if handler != nil {
		return handler, nil
	}

	return nil, fmt.Errorf("no handler for: '%s' resource found", resourceType)
}

func (p *ExecutionPlan) SetContext(c *context.Context) {
	p.ctx = c
}

func (p *ExecutionPlan) Plan(resources *map[string]config.Resource, sess *session.Session) error {
	// TODO: Add logic for dependency tree
	for _, resource := range *resources {
		handler, err := handlerFactory(resource.Type, resource.Options, sess)
		if err != nil {
			return err
		}

		p.Steps = append(p.Steps, ExecutionStep{
			Type:    resource.Type,
			handler: handler,
		})
	}

	return nil
}

func (p *ExecutionPlan) Exec(err *errors.ErrorsBag) error {
	if len(p.Steps) < 0 {
		return fmt.Errorf("%v execution steps provided, skipping", len(p.Steps))
	}

	var wg sync.WaitGroup

	wg.Add(len(p.Steps))

	for _, step := range p.Steps {
		go func(step ExecutionStep) {
			defer wg.Done()
			err := step.handler.Run()
			if err != nil {
				fmt.Printf("Error executing step for: '%s', Error: %v", step.Type, err.Error())
			}
		}(step)
	}

	wg.Wait()

	return nil
}
