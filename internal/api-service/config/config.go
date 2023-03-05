package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/abhirajranjan/gochat/pkg/constants"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	configfile string
)

func init() {
	flag.StringVar(&configfile, "conf", "", "gochat config file path")
}

func LoadConf() (*Config, error) {
	if configfile == "" {
		configFilePath := os.Getenv(constants.CONFIG_FILE_ENV)
		if configFilePath != "" {
			configfile = configFilePath
		} else {
			wd, err := os.Getwd()
			if err != nil {
				errors.Wrap(err, "os.Getwd")
			}
			configfile = fmt.Sprintf("%s/config.yaml", wd)
		}
	}

	cfg := &Config{}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configfile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(&cfg, mapstructureHooks()...); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	if cfg.Logger.AppName == "" {
		cfg.Logger.AppName = cfg.AppName
	}

	return cfg, nil
}

func mapstructureHooks() []viper.DecoderConfigOption {
	hooks := []viper.DecoderConfigOption{
		stringToTimeDurationHookFunc(),
	}
	return hooks
}

// TODO: add plugin architecture to add hooks
func stringToTimeDurationHookFunc() viper.DecoderConfigOption {
	return func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		)
	}
}
