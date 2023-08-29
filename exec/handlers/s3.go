package handlers

import (
	"fmt"
	"localstack-quickstart/config"
)

type S3Handler struct {
	Options config.S3Options
}

func (s S3Handler) Run() error {
	fmt.Printf("Executing S3, name: %s\n", s.Options.Name)

	return nil
}
