package config

import (
	"os"

	"gmvr.pw/boombox-web-runner/pkg/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Entrypoints EntrypointsConfig `yaml:"entrypoints"`
	Modules     []model.Module    `yaml:"modules"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{}

	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		err = yaml.NewDecoder(file).Decode(&cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
