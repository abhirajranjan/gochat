package payloadManager

import (
	"sync"

	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/pkg/errors"
)

const (
	// key to add to payload (map[string]interface{} type) for version identification
	PAYLOAD_VERSION_TAG = "version"
)

// manager errors
var (
	ErrPluginAlreadyExist = errors.New("plugin already exist")
	ErrNoParserFound      = errors.New("no parser found")
	ErrParser             = errors.New("parser failed to process request")
	ErrBadPayload         = errors.New("bad payload")
)

// parser errors
// registered parsers need to return same error string for consistency
var (
	ErrNoDataInPayload           = errors.New("no data in payload")
	ErrFailedToDecodePayloadData = errors.New("failed to decode payload")
)

type IPayloadData interface {
	Version() int64
	Get(string) (interface{}, bool)
}

type ILoginResponse interface {
	GetMap() (map[string]interface{}, error)
}

type IParser interface {
	// should return supported version
	//
	// must not collide with already existing version numbers
	SupportsVersion() int64

	// encode should return map after encoding the data.GetMap() into it
	//
	// should return error returned by data.GetMap() if failed.
	Encode(data interface {
		GetMap() (map[string]interface{}, error)
	}, sessionID string) (map[string]interface{}, error)

	// decode retrives data set by Encode function
	//
	// should return error ErrNoDataInPayload if payload doesnt contain data as set by parser logic
	// return ErrFailedToDecodePayloadData if data key exists in payload but is malformatted
	Decode(map[string]interface{}) (interface {
		Version() int64
		Get(string) (interface{}, bool)
	}, error)
}

type Manager struct {
	logger         logger.ILogger
	mu             sync.RWMutex
	parsers        map[int64]IParser
	latestVersion  int64
	minimumVersion int64
}

func NewManager(logger logger.ILogger) *Manager {
	return &Manager{logger: logger, parsers: map[int64]IParser{}}
}

// register new parser of type IParser
//
// returns ErrPluginAlreadyExist if parser with same version already exists
// else returns nil
func (m *Manager) RegisterParser(parser IParser) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.parsers[parser.SupportsVersion()]; ok {
		m.logger.Warnf("parser having version %d already registered", parser.SupportsVersion())
		return ErrPluginAlreadyExist
	}

	m.logger.Infof("added parser version %d", parser.SupportsVersion())
	m.parsers[parser.SupportsVersion()] = parser

	if parser.SupportsVersion() > m.latestVersion {
		m.latestVersion = parser.SupportsVersion()
	}
	return nil
}

// fetch parser from map
func (m *Manager) getParser(version int64) IParser {
	m.mu.RLock()
	defer m.mu.RUnlock()

	parser, ok := m.parsers[version]
	if !ok {
		return nil
	}

	return parser
}

// set minimum version for the parser
//
// minimum version should be long term and should be able to pass almost every request
//
// minimum version is used as a last chance to decode and ecode in case of any error is raised by latest parser
func (m *Manager) SetMinimumVersion(minimiumVersion int64) error {
	if parser := m.getParser(minimiumVersion); parser != nil {
		m.logger.Infof("added minimum version parser: %d", minimiumVersion)
		m.minimumVersion = minimiumVersion
		return nil
	}
	m.logger.Warnf("given minimum version parser (version: %d) is not registed", minimiumVersion)
	return ErrNoParserFound
}

func (m *Manager) GetMinimumVersion() int64 {
	return m.minimumVersion
}

// encode the data into out map with respective version parser
//
// returns ErrNoParserFound if no parser found.
//
// else if err != nil then err is returned by data.GetMap() call
func (m *Manager) encodeWithVersion(data ILoginResponse, sessionID string, out map[string]interface{}, version int64) error {
	parser := m.getParser(version)
	if parser == nil {
		m.logger.Debugf("encode parser version: %d  not found", version)
		return ErrNoParserFound
	}

	o, err := parser.Encode(data, sessionID)
	if err != nil {
		m.logger.Debugf("parser (version %d) failed to encode with err: %s", version, err)
		return err
	}

	for key, val := range o {
		out[key] = val
	}
	out[PAYLOAD_VERSION_TAG] = parser.SupportsVersion()
	m.logger.Debugf("parser (version %d) encode result: %#v", out)
	return nil
}

// encode encode data into out with latest parser. fallback to minimum version if parsing with latest parser failed.
//
// returns ErrNoParserFound if no parser exists.
//
// if err != nil then err is be returned by data.GetMap() call
func (m *Manager) Encode(data interface {
	GetMap() (map[string]interface{}, error)
}, sessionID string, out map[string]interface{}) error {
	err := m.encodeWithVersion(data, sessionID, out, m.latestVersion)
	if err != nil {
		m.logger.Debugf("latest parser failed to encode data")
		err := m.encodeWithVersion(data, sessionID, out, m.GetMinimumVersion())
		return err
	}
	return nil
}

// decode data with version parser
//
// returns ErrNoParserFound if no parser found.
//
// returns ErrNoPayloadData if payload data is not in expected format.
//
// returns ErrFailedToDecodePayloadData if payload data exists but in wrong format
func (m *Manager) decodeWithVersion(data map[string]interface{}, version int64) (IPayloadData, error) {
	parser := m.getParser(version)
	if parser == nil {
		m.logger.Warnf("parser (version %d) not found for decoding", version)
		return nil, ErrNoParserFound
	}
	parsedData, err := parser.Decode(data)
	if err != nil {
		m.logger.Debugf("parser (version %d) fails to decode: %s", err)
	} else {
		m.logger.Debugf("parser (version %d) decode response: %#v", parsedData)
	}
	return parsedData, err
}

// decode the data into payload to get user info in jwt
//
// returns ErrBadPayload if parser doesnt have version tag that stores version which encode payload
//
// returns ErrNoParserFound if no parser found.
//
// returns ErrNoPayloadData if payload data is not in expected format.
//
// returns ErrFailedToDecodePayloadData if payload data exists but in wrong format
func (m *Manager) Decode(data map[string]interface{}) (interface {
	Version() int64
	Get(string) (interface{}, bool)
}, error) {

	ver := getVersion(data)
	if ver < 0 {
		handleInvalidVersion(ver, m.logger)
		return nil, ErrBadPayload
	}

	return m.decodeWithVersion(data, ver)
}

// returns the payload version
//
// returns -1 if PAYLOAD_VERSION_TAG exists but is not int,
// -2 if it does not exist
func getVersion(data map[string]interface{}) int64 {
	if ver, ok := data[PAYLOAD_VERSION_TAG]; ok {
		if version, ok := ver.(float64); ok {
			return int64(version)
		}
		return -1
	}
	return -2
}

// handle logging of errors in getVersion calls
func handleInvalidVersion(ver int64, logger logger.ILogger) {
	if ver == -1 {
		logger.Debugf("invalid version provided in payload (%T): %v", ver, ver)
	}
	if ver == -2 {
		logger.Debug("no version provided in payload")
	}
}
