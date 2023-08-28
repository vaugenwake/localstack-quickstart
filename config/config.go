package config

import (
	"io/ioutil"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Connection struct {
	Protocol string `yaml:"protocol"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

type Config struct {
	Connection Connection `yaml:"connection"`
}

func (c *Config) GetEndpoint() string {
	return c.Connection.Protocol + "://" + c.Connection.Endpoint + ":" + strconv.Itoa(c.Connection.Port)
}

func ParseConfigFile(file string) (*Config, error) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.UnmarshalStrict(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
