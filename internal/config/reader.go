package config

import (
	"github.com/kelseyhightower/envconfig"
)

func Read(prefix string) (*Config, error) {
	var cfg Config

	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
