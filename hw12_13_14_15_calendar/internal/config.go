package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Server ServerConf
	Logger LoggerConf
	DB     DBConf
}

type ServerConf struct {
	BindAddress string        `config:"bind_address,require"`
	ReadTimeout time.Duration `config:"read_timeout"`
	WriteTimeot time.Duration `config:"write_timeout"`
}

type LoggerConf struct {
	Level string `config:"level"`
}

type DBConf struct {
	DSN string `config:"DSN,require"`
}

func NewConfig(configFile string) (*Config, error) {
	f, _ := os.Open(configFile)
	content, _ := ioutil.ReadAll(f)
	fmt.Println(string(content))

	cfg := Config{
		Server: ServerConf{},
		DB:     DBConf{},
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
