package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getStubPath(stub_name string) string {
	return filepath.Join("../../testdata", stub_name)
}

func TestCanParseConfig(t *testing.T) {
	assert := assert.New(t)

	validResources := map[string]Resource{
		"my-bucket": {
			Type: S3,
			Options: S3Options{
				Name: "test-bucket",
			},
		},
		"my-queue": {
			Type: SQS,
			Options: SQSOptions{
				Name:       "test-queue",
				DeadLetter: false,
			},
		},
	}

	var tests = []struct {
		name               string
		stub               string
		expectError        bool
		expectedConnection Connection
		expectedResources  map[string]Resource
	}{
		{
			name:        "valid config",
			stub:        "valid_config.yml",
			expectError: false,
			expectedConnection: Connection{
				Protocol: "http",
				Endpoint: "localstack-test",
				Port:     4566,
			},
			expectedResources: validResources,
		},
		{
			name:               "invalid config",
			stub:               "invalid_config.yml",
			expectError:        true,
			expectedConnection: Connection{},
			expectedResources:  map[string]Resource{},
		},
		{
			name:               "invalid resource options",
			stub:               "invalid_resource_options.yml",
			expectError:        true,
			expectedConnection: Connection{},
			expectedResources:  map[string]Resource{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := getStubPath(test.stub)

			result, err := ParseConfigFile(stub)

			if test.expectError && err == nil {
				assert.Error(err)
				return
			}

			if !test.expectError {
				assert.Equal(test.expectedConnection, result.Connection)
				assert.Equal(test.expectedResources, result.Resources)
			}

		})
	}
}

func TestCanBuildEndpoint(t *testing.T) {
	assert := assert.New(t)

	stub := getStubPath("valid_config.yml")

	result, err := ParseConfigFile(stub)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.Equal("http://localstack-test:4566", result.Connection.GetEndpoint())
}

func TestCanMarshalOptionsTypeForResource(t *testing.T) {
	assert := assert.New(t)

	stub := getStubPath("valid_config.yml")

	result, err := ParseConfigFile(stub)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.IsType(S3Options{}, result.Resources["my-bucket"].Options)
}
