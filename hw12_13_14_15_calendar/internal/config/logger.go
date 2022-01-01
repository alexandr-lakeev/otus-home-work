package config

type LoggerConf struct {
	Env   string `config:"ENV"`
	Level string `config:"level"`
}
