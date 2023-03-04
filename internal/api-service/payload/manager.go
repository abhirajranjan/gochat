package payload

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/pkg/logger"
)

var (
	ErrPluginAlreadyExist = errors.New("plugin with supported version already exist")
	ErrNoParserFound      = errors.New("no parser found for given version")
	ErrParser             = errors.New("parser failed to process request")
)

type Manager struct {
	logger        logger.ILogger
	mu            sync.RWMutex
	parsers       map[int64]model.IParser
	latestVersion int64

	//config

	// if inplace is set to true then all parser functions
	// will do inplace parsing and conversely if set to false
	// then parsers will return type interface{} having modified value
	inplace bool
}

func NewManager(logger logger.ILogger, inplace bool) model.IPayLoadManager {
	return &Manager{logger: logger, inplace: inplace}
}

func (m *Manager) RegisterParser(parser model.IParser) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.parsers[parser.SupportsVersion()]; ok {
		return ErrPluginAlreadyExist
	}
	m.parsers[parser.SupportsVersion()] = parser

	if parser.SupportsVersion() > m.latestVersion {
		m.latestVersion = parser.SupportsVersion()
	}
	return nil
}

func (m *Manager) getParser(version int64) model.IParser {
	m.mu.RLock()
	defer m.mu.Unlock()

	parser, ok := m.parsers[version]
	if !ok {
		return nil
	}

	return parser
}

func (m *Manager) Encode(data map[string]interface{}, version int64) (map[string]interface{}, error) {
	parser := m.getParser(version)
	if parser != nil {
		return parser.Encode(data, m.inplace)
	}
	return nil, ErrNoParserFound
}

func (m *Manager) AddPayload(data map[string]interface{}) (map[string]interface{}, error) {
	return m.Encode(data, m.latestVersion)
}

func (m *Manager) Decode(data map[string]interface{}, version int64) (model.IPayloadData, error) {
	parser := m.getParser(version)
	if parser == nil {
		return nil, ErrNoParserFound
	}
	return parser.Decode(data)
}

func (m *Manager) VerifyUser(data model.IPayloadData) bool {
	parser := m.getParser(data.Version())
	if parser == nil {
		return false
	}
	return parser.VerifyUser(data)
}

func (m *Manager) LogoutUser(data map[string]interface{}, version int64) bool {
	parser := m.getParser(version)
	if parser == nil {
		return false
	}
	err := parser.LogoutUser(data)
	if err != nil {
		err = errors.Wrap(err, "ParserManager.LogoutUser")
		m.logger.Error(err)
		return false
	}
	return true
}
