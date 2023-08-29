package validation

import (
	"fmt"
	"localstack-quickstart/config"

	"slices"
)

type ValidationError struct {
	Field   string
	Message string
}

type ValidationReport struct {
	errors []ValidationError
}

func (r *ValidationReport) AddError(field string, message string) {
	r.errors = append(r.errors, ValidationError{Field: field, Message: message})
}

func ValidateResourceTypeField(resourceType config.ResourceType) error {
	validTypes := []config.ResourceType{
		config.S3,
		config.SQS,
	}

	if !slices.Contains(validTypes, resourceType) {
		return fmt.Errorf("resources.type not supported: %v", resourceType)
	}

	return nil
}

func (r *ValidationReport) HasErrors() bool {
	return len(r.errors) > 0
}

func (r *ValidationReport) AllErrors() []ValidationError {
	return r.errors
}

func (r *ValidationReport) GetErrorAtIndex(index int) (*ValidationError, error) {
	if len(r.errors)-1 < index {
		return nil, fmt.Errorf("GetErrorAtIndex(), out of bounds, index: %v, slice length: %v", index, len(r.errors))
	}

	return &r.errors[index], nil
}

func (r *ValidationReport) ValidateResources(resources *config.Resourses) {

	for resourceName, resource := range *resources {
		err := ValidateResourceTypeField(resource.Type)
		if err != nil {
			r.AddError(fmt.Sprintf("resources.%s.type", resourceName), err.Error())
			continue
		}
	}
}
