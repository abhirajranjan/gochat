package payloadParserV1

import "errors"

const (
	VERSION             int64  = 1
	DATA_PAYLOAD_TAG    string = "data"
	SESSION_PAYLOAD_TAG string = "session"
	COMMAND_SESSION_TAG string = "__session__"
)

var (
	ErrNoDataInPayload           = errors.New("no data in payload")
	ErrFailedToDecodePayloadData = errors.New("failed to decode payload data")
	ErrNoSessionInPayload        = errors.New("no sessionID in payload")
)

type Payload map[string]interface{}

func (p Payload) setSessionID(a any) {
	p[COMMAND_SESSION_TAG] = a
}

func (p Payload) Version() int64 {
	return VERSION
}

func (p Payload) GetSessionID() interface{} {
	return p[COMMAND_SESSION_TAG]
}

func (p Payload) Get(key string) (interface{}, bool) {
	switch key {
	case "sessionID":
		return p[COMMAND_SESSION_TAG], true
	default:
		value, ok := p[key]
		return value, ok
	}
}

type v1Parser struct{}

func NewV1Parser() *v1Parser {
	return &v1Parser{}
}

func (p *v1Parser) SupportsVersion() int64 {
	return VERSION
}

func (p *v1Parser) Encode(data interface {
	GetMap() (map[string]interface{}, error)
}, sessionID string) (map[string]interface{}, error) {

	out := make(map[string]interface{})
	mapData, err := data.GetMap()
	if err != nil {
		return nil, err
	}
	out[SESSION_PAYLOAD_TAG] = sessionID
	out[DATA_PAYLOAD_TAG] = mapData
	return out, nil
}

func (p *v1Parser) Decode(data map[string]interface{}) (interface {
	Version() int64
	Get(string) (interface{}, bool)
	GetSessionID() interface{}
}, error) {

	mapData, ok := data[DATA_PAYLOAD_TAG]
	if !ok {
		return nil, ErrNoDataInPayload
	}

	mappedData, ok := mapData.(map[string]interface{})
	if !ok {
		return nil, ErrFailedToDecodePayloadData
	}

	mapSession, ok := data[SESSION_PAYLOAD_TAG]
	if !ok {
		return nil, ErrNoSessionInPayload
	}

	mappedSession, ok := mapSession.(map[string]interface{})
	if !ok {
		return nil, ErrFailedToDecodePayloadData
	}

	var payload Payload = mappedData
	payload.setSessionID(mappedSession)
	return &payload, nil
}
