package exec

import (
	"fmt"
	"localstack-quickstart/config"
	"localstack-quickstart/exec/handlers"
)

type Handler interface {
	Run() error
}

type ExecutionStep struct {
	Type    config.ResourceType
	handler Handler
}

type ExecutionPlan struct {
	Steps []ExecutionStep
}

func handlerFactory(resourceType config.ResourceType, options interface{}) (Handler, error) {
	var handler Handler

	switch resourceType {
	case config.S3:
		if s3Opts, ok := options.(config.S3Options); ok {
			handler = &handlers.S3Handler{Options: s3Opts}
		}
	case config.SQS:
		if sqsOpts, ok := options.(config.SQSOptions); ok {
			handler = &handlers.SQSHandler{
				Options: sqsOpts,
			}
		}
	}

	if handler != nil {
		return handler, nil
	}

	return nil, fmt.Errorf("No handler for: '%s' resource found", resourceType)
}

func (p *ExecutionPlan) Plan(resources *map[string]config.Resource) error {
	// TODO: Add logic for dependency tree
	for _, resource := range *resources {
		handler, err := handlerFactory(resource.Type, resource.Options)
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

func (p *ExecutionPlan) Exec() error {
	if len(p.Steps) < 0 {
		return fmt.Errorf("%v execution steps provided, skipping", len(p.Steps))
	}

	for _, step := range p.Steps {
		err := step.handler.Run()
		if err != nil {
			fmt.Printf("Error executing step for: '%s', Error: %v", step.Type, err.Error())
		}
	}

	return nil
}
