package config

type LoggerConfig struct {
	Name    string `config:"APP_NAME"`
	Level   string `config:"LOGGER_LEVEL"`
	Isdev   bool   `config:"LOGGER_ISDEV"`
	Encoder string `config:"LOGGER_ENCODER"`
}
