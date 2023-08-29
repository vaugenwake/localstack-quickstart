package handlers

import (
	"fmt"
	"localstack-quickstart/config"
)

type SQSHandler struct {
	Options config.SQSOptions
}

func (S SQSHandler) Run() error {
	fmt.Printf("Executing SQS, name: %s\n", S.Options.Name)

	return nil
}
