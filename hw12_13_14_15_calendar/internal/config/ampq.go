package config

type AmpqConf struct {
	URL          string `config:"AMPQ_URL"`
	ExchangeType string `toml:"exchange_type"`
	ExchangeName string `toml:"exchange_name"`
	QueueName    string `toml:"queue_name"`
}
