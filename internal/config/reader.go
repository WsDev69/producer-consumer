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

func ReadProducer(prefix string) (*Producer, error) {
	var cfg Producer

	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ReadConsumer(prefix string) (*Consumer, error) {
	var cfg Consumer

	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
