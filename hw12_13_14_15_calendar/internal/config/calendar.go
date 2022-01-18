package config

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

const (
	MemoryStorage = "memory"
	SQLStorage    = "sql"
)

type (
	Config struct {
		Server  ServerConf  `toml:"server"`
		Grpc    GrpcConf    `toml:"grpc"`
		Logger  LoggerConf  `toml:"logger"`
		Storage StorageConf `toml:"storage"`
	}

	ServerConf struct {
		BindAddress  string        `config:"http_bind_address,require"`
		ReadTimeout  time.Duration `config:"http_read_timeout"`
		WriteTimeout time.Duration `config:"http_write_timeout"`
		IdleTimeout  time.Duration `config:"http_idle_timeout"`
	}

	GrpcConf struct {
		ListenAddress string `config:"grpc_listen_address,require"`
	}
)

func NewConfig(configFile string) (*Config, error) {
	cfg := Config{
		Server: ServerConf{
			BindAddress: ":8080",
		},
		Grpc: GrpcConf{
			ListenAddress: ":50051",
		},
		Storage: StorageConf{
			Type: MemoryStorage,
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
