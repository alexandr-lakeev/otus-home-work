package internal

import "os"

type Config struct {
	Logger LoggerConf
	DB     DBConf
}

type LoggerConf struct {
	Level string
}

type DBConf struct {
	DSN string
}

func NewConfig() Config {
	return Config{
		Logger: LoggerConf{
			Level: "INFO", // TODO from file
		},
		DB: DBConf{
			DSN: os.Getenv("DSN"),
		},
	}
}

// TODO
