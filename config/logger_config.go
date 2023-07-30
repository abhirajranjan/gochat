package config

type LoggerConfig struct {
	AppName string `mapstructure:"appName"`
	Level   string `mapstructure:"level"`
	Dev     bool   `mapstructure:"devMode"`
	Encoder string `mapstructure:"encoder"`
}
