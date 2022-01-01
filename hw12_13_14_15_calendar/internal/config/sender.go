package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type (
	SenderConfig struct {
		Storage StorageConf
		Ampq    AmpqConf   `toml:"rabbitmq"`
		Logger  LoggerConf `toml:"logger"`
	}
)

func NewSenderConfig(configFile string) (*SenderConfig, error) {
	cfg := SenderConfig{}

	err := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(configFile),
	).Load(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
