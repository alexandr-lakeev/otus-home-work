package config

type StorageConf struct {
	Type string `config:"type"`
	DSN  string `config:"DB_DSN"`
}
