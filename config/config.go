package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type ResourceType string

var (
	S3  ResourceType = "s3"
	SQS ResourceType = "sqs"
)

type Connection struct {
	Protocol string `yaml:"protocol"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

type S3Options struct {
	Name string `yaml:"name"`
}

type SQSOptions struct {
	Name       string `yaml:"name"`
	DeadLetter bool   `yaml:"deadLetter"`
}

type Resource struct {
	Type    ResourceType `yaml:"type"`
	Options interface{}  `yaml:"options"`
}

type Config struct {
	Connection Connection          `yaml:"connection"`
	Resources  map[string]Resource `yaml:"resources"`
}

func (r *Resource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type resourceAlias Resource
	var ra resourceAlias

	if err := unmarshal(&ra); err != nil {
		return err
	}

	r.Type = ra.Type

	switch ra.Type {
	case "s3":
		var c struct {
			Options S3Options `yaml:"options"`
		}
		err := unmarshal(&c)
		r.Options = c.Options
		return err
	case "sqs":
		var c struct {
			Options SQSOptions `yaml:"options"`
		}

		if err := unmarshal(&c); err != nil {
			return err
		}

		r.Options = c.Options
	default:
		return fmt.Errorf("Unknown resource type: %v", ra.Type)
	}

	return nil
}

func (c *Config) GetEndpoint() string {
	return c.Connection.Protocol + "://" + c.Connection.Endpoint + ":" + strconv.Itoa(c.Connection.Port)
}

func ParseConfigFile(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
