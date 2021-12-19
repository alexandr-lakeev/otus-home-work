package config

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

const STORAGE_MEMORY = "memory"
const STORAGE_SQL = "sql"

type Config struct {
	Server  ServerConf
	Logger  LoggerConf
	Storage StorageConf
}

type ServerConf struct {
	BindAddress  string        `config:"bind_address,require"`
	ReadTimeout  time.Duration `config:"read_timeout"`
	WriteTimeout time.Duration `config:"write_timeout"`
	IdleTimeout  time.Duration `config:"idle_timeout"`
}

type LoggerConf struct {
	Env   string `config:"ENV"`
	Level string `config:"level"`
}

type StorageConf struct {
	Type string `config:"type"`
	DSN  string `config:"DSN"`
}

func NewConfig(configFile string) (*Config, error) {
	cfg := Config{
		Server: ServerConf{
			BindAddress: ":8080",
		},
		Storage: StorageConf{
			Type: STORAGE_MEMORY,
		},
		Logger: LoggerConf{
			Level: "INFO",
		},
	}

	err := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(configFile),
	).Load(context.Background(), &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
