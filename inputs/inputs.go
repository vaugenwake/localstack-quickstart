package inputs

import (
	"flag"
)

type Inputs struct {
	ConfigFile string
}

func ParseInputFlags() (*Inputs, error) {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "", "Path of config file")

	flag.Parse()

	return &Inputs{
		ConfigFile: configFilePath,
	}, nil
}
