package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type (
	SchedulerConfig struct {
		Storage StorageConf
		Ampq    AmpqConf `toml:"rabbitmq"`
	}

	AmpqConf struct {
		URL          string `config:"AMPQ_URL"`
		ExchangeType string `toml:"exchange_type"`
		ExchangeName string `toml:"exchange_name"`
		QueueName    string `toml:"queue_name"`
	}
)

func NewSchedulerConfig(configFile string) (*SchedulerConfig, error) {
	cfg := SchedulerConfig{}

	err := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(configFile),
	).Load(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
