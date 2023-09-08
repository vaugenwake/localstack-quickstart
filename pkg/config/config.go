package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
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
	DeadLetter bool   `yaml:"dead_letter"`
}

type Resource struct {
	Type    ResourceType `yaml:"type"`
	Options interface{}  `yaml:"options"`
}

type Config struct {
	Connection Connection          `yaml:"connection"`
	Resources  map[string]Resource `yaml:"resources"`
}

func (r *Resource) UnmarshalYAML(value *yaml.Node) error {
	type resourceAlias Resource
	var ra resourceAlias

	if err := value.Decode(&ra); err != nil {
		return err
	}

	r.Type = ra.Type

	switch ra.Type {
	case "s3":
		var c struct {
			Options S3Options `yaml:"options"`
		}

		if err := value.Decode(&c); err != nil {
			return err
		}

		if c.Options == (S3Options{}) {
			return fmt.Errorf("options not provided for: %v", r.Type)
		}

		r.Options = c.Options
	case "sqs":
		var c struct {
			Options SQSOptions `yaml:"options"`
		}

		if err := value.Decode(&c); err != nil {
			return err
		}

		if c.Options == (SQSOptions{}) {
			return fmt.Errorf("options not provided for: %v", r.Type)
		}

		r.Options = c.Options
	default:
		return fmt.Errorf("Unknown resource type: %v", ra.Type)
	}

	return nil
}

func (c *Connection) GetEndpoint() string {
	return c.Protocol + "://" + c.Endpoint + ":" + strconv.Itoa(c.Port)
}

func ParseConfigFile(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
