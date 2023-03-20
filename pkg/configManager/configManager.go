package configManager

import (
	"errors"
	"log"
)

var (
	ErrConfigParserAlreadyExist = errors.New("config parser having indential parsing type exist")
)

// parser interface to register to manager
type ConfigParser interface {
	LoadConfig(interface{}) error
	GetParsingType() string
}

// function returns error if cfg is not valid else other middlefunctions get called
type Middlewarefunc func(cfg any) error

type ConfigManager[T any] struct {
	mappingorder []ConfigParser

	// functions that acts as middleware for adding various middleware methods
	// like logging, domain checks ...
	// get called serially for each parseable parser
	// returns error if cfg could not proceed, error will be shown as output
	Middlewarefunc []Middlewarefunc
}

// generate new ConfigManager that decodes data into generic type T by LoadConfig
func NewConfigManager[T any]() *ConfigManager[T] {
	confManager := ConfigManager[T]{
		mappingorder: make([]ConfigParser, 0),
	}
	return &confManager
}

// register a new parser of type ConfigParser
//
// new parser will be called in order in which it is registered.
// register parser in priority order
func (m *ConfigManager[T]) RegisterConfigParser(parser ConfigParser) {
	m.mappingorder = append(m.mappingorder, parser)
}

// add middleware functions to config manager
func (m *ConfigManager[T]) AddMiddleware(f ...Middlewarefunc) {
	m.Middlewarefunc = append(m.Middlewarefunc, f...)
}

// loads config into object of generic type T
//
// calls parser serially in order in which they are registered
// if parser fails to process it will log error raised by parser.LoadConfig
// and move to the next parser.
//
// if middleware functions are specified, they will get called
// after parser successfully parsed into generic type T, this generic type is
// passed to middleware functions. if all middleware function returns nil (no error)
// then this function returns returns parsed object pointer
func (m *ConfigManager[T]) LoadConfig() *T {
	for _, parser := range m.mappingorder {
		cfg := new(T)
		err := parser.LoadConfig(cfg)
		if err == nil || cfg != nil {
			for _, middlewarefunc := range m.Middlewarefunc {
				if err := middlewarefunc(cfg); err != nil {
					logMiddlewarefuncErr(err)
					return nil
				}
			}
			return cfg
		}
		logParserError(parser, err)
	}
	logNoParserPassed()
	return nil
}

// logging function if middleware function fails
func logMiddlewarefuncErr(err error) {
	log.Fatalf("middleware function failed with err: %s", err.Error())
}

// logging function if all registered parser fail to parse
func logNoParserPassed() {
	log.Fatal("no parser found to load config")
}

// logging function if parser fails
//
// logs error and then pass on to next parser in list
func logParserError(parser ConfigParser, err error) {
	// TODO: add verbose level
	log.Printf("parser of type %s failed with err: %s\n", parser.GetParsingType(), err.Error())
}
