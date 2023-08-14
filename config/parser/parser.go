package parser

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load(cfg any, configfile string) error {
	parser := viper.New()

	replacer := strings.NewReplacer(".", "_")
	parser.SetEnvKeyReplacer(replacer)
	parser.AutomaticEnv()

	if configfile != "" {
		parser.SetConfigFile(configfile)

		if err := parser.ReadInConfig(); err != nil {
			return fmt.Errorf("viper.ReadInConfig: %w", err)
		}
	}
	parser.Get("sql.dsn")
	if err := parser.Unmarshal(&cfg, mapstructureHooks()...); err != nil {
		return fmt.Errorf("viper.Marshal: %w", err)
	}

	return nil
}

func mapstructureHooks() []viper.DecoderConfigOption {
	hooks := []viper.DecoderConfigOption{
		stringToTimeDurationHookFunc(),
	}
	return hooks
}

func stringToTimeDurationHookFunc() viper.DecoderConfigOption {
	return func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = mapstructure.StringToTimeDurationHookFunc()
	}
}
