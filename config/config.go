package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Connection struct {
	Protocol string `yaml:"protocol"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

type Resource struct {
	Type    string      `yaml:"type"`
	Options interface{} `yaml:"options"`
}

type Config struct {
	Connection Connection          `yaml:"connection"`
	Resources  map[string]Resource `yaml:"resources"`
}

func (c *Config) GetEndpoint() string {
	return c.Connection.Protocol + "://" + c.Connection.Endpoint + ":" + strconv.Itoa(c.Connection.Port)
}

func ParseConfigFile(file string) (*Config, error) {
	configFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.UnmarshalStrict(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
