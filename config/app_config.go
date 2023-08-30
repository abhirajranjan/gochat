package config

type AppConfig struct {
	Addr string `config:"APP_ADDR"`
	Port string `config:"APP_PORT"`
}
