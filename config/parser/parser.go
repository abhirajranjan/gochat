package parser

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	ErrNoConfigFile      = errors.New("no yaml config file provided")
	ErrInvalidConfigFile = errors.New("config file not found")
)

func Load(cfg any, configfile string) error {
	parser := viper.New()
	parser.SetConfigFile(configfile)

	if err := parser.ReadInConfig(); err != nil {
		return fmt.Errorf("viper.ReadIn: %w", err)
	}

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