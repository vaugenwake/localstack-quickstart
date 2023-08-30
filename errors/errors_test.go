package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	errorsBag := &ErrorsBag{}

	errorsBag.Add("error", "This is a test error")

	assert.Equal(1, len(errorsBag.All()))
	assert.Equal(Error{"error", "This is a test error"}, errorsBag.Get(0))
}

func TestAny(t *testing.T) {
	assert := assert.New(t)

	errorsBag := &ErrorsBag{}

	errorsBag.Add("error", "This is a test error")
	assert.True(errorsBag.Any())

	errorsBagEmpty := &ErrorsBag{}
	assert.False(errorsBagEmpty.Any())
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	errorsBag := &ErrorsBag{}

	errorsBag.Add("error", "Error 1")
	errorsBag.Add("error", "Error 2")

	assert.Equal(2, len(errorsBag.All()))
	assert.Equal(Error{"error", "Error 2"}, errorsBag.Get(1))
}

func TestAll(t *testing.T) {
	assert := assert.New(t)

	expectedResult := []Error{
		{"log", "entry 1"},
		{"log", "entry 2"},
	}

	errorsBag := &ErrorsBag{}

	errorsBag.Add("log", "entry 1")
	errorsBag.Add("log", "entry 2")

	assert.Equal(expectedResult, errorsBag.All())
}
