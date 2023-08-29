package validation

import (
	"localstack-quickstart/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanAddValidationError(t *testing.T) {
	assert := assert.New(t)

	validator := &ValidationReport{}

	validator.AddError("test", "This is an error")

	expected := &ValidationError{
		Field:   "test",
		Message: "This is an error",
	}

	assert.Equal(1, len(validator.AllErrors()))

	addedError, _ := validator.GetErrorAtIndex(0)

	assert.Equal(expected, addedError)
}

func TestValidateResourceType(t *testing.T) {
	assert := assert.New(t)

	result := ValidateResourceTypeField("s3")
	assert.Nil(result)

	result = ValidateResourceTypeField("apigateway")
	assert.Error(result)
}

func TestHasErrors(t *testing.T) {
	assert := assert.New(t)

	validator := &ValidationReport{}

	validator.AddError("type", "error 1")
	validator.AddError("type", "error 2")

	result := validator.HasErrors()

	assert.True(result)
}

func TestAllErrors(t *testing.T) {
	assert := assert.New(t)

	validator := &ValidationReport{}

	validator.AddError("type", "error 1")
	validator.AddError("type", "error 2")

	result := validator.AllErrors()

	assert.Equal(2, len(result))
	assert.Equal(ValidationError{Field: "type", Message: "error 1"}, result[0])
}

func TestGetErrorAtIndex(t *testing.T) {
	assert := assert.New(t)

	validator := &ValidationReport{}

	validator.AddError("type", "error 1")
	validator.AddError("type", "error 2")

	result, err := validator.GetErrorAtIndex(0)
	if err != nil {
		t.Errorf("Unexpected error performing: GetErrorAtIndex(n), %v", err)
	}

	assert.Equal(&ValidationError{Field: "type", Message: "error 1"}, result)

	result, err = validator.GetErrorAtIndex(2)
	assert.Error(err)
	assert.Nil(result)
}

func TestValidateSliceOfResources(t *testing.T) {
	assert := assert.New(t)

	validator := &ValidationReport{}

	resources := &config.Resourses{
		"bucket": {
			Type: "s3",
			Options: map[string]interface{}{
				"name": "bucket",
			},
		},
		"api": {
			Type: "apigateway",
			Options: map[string]interface{}{
				"name": "api",
			},
		},
		"queue": {
			Type: "sqs",
			Options: map[string]interface{}{
				"name": "queue",
			},
		},
	}

	validator.ValidateResources(resources)

	assert.Equal(1, len(validator.AllErrors()))
}
