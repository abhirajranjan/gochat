package yamlParser

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	ErrNoConfigFile      = errors.New("no yaml config file provided")
	ErrInvalidConfigFile = errors.New("config file not found")
)

type yamlParser struct {
	configFile  string
	viperParser *viper.Viper
}

func NewYamlParser(configfile string) *yamlParser {
	parser := &yamlParser{
		configFile:  configfile,
		viperParser: viper.New(),
	}
	parser.viperParser.SetConfigType("yaml")
	return parser
}

func (parser *yamlParser) Load(cfg any) error {

	if parser.configFile == "" {
		return ErrNoConfigFile
	}

	file, err := os.OpenFile(parser.configFile, os.O_RDONLY, 0)
	if errors.Is(err, os.ErrNotExist) {
		return ErrInvalidConfigFile
	}

	if err := parser.viperParser.ReadConfig(file); err != nil {
		return fmt.Errorf("viper.ReadIn: %w", err)
	}

	if err := viper.Unmarshal(&cfg, mapstructureHooks()...); err != nil {
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
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		)
	}
}
