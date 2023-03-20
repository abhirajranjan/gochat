package yamlParser

import (
	"flag"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const TYPE = "yaml"

func init() {
	flag.StringVar(&configfile, "yaml", "", "yaml config file")
}

var (
	ErrNoConfigFile      = errors.New("no yaml config file provided")
	ErrInvalidConfigFile = errors.New("config file not found")
	configfile           string
)

type yamlParser struct {
	configFile string
}

func NewYamlParser(configfile string) *yamlParser {
	return &yamlParser{configFile: configfile}
}

func (parser *yamlParser) GetParsingType() string {
	return TYPE
}

func (parser *yamlParser) LoadConfig(cfg interface{}) error {
	if parser.configFile == "" {
		return ErrNoConfigFile
	}

	file, err := os.OpenFile(parser.configFile, os.O_RDONLY, 0000)
	if errors.Is(err, os.ErrNotExist) {
		return ErrInvalidConfigFile
	}

	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(file); err != nil {
		return errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(&cfg, mapstructureHooks()...); err != nil {
		return errors.Wrap(err, "viper.Unmarshal")
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
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		)
	}
}
