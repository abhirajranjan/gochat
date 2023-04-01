package mockParser

import "errors"

const (
	VERSION          int64  = 1
	DATA_PAYLOAD_TAG string = "data"
)

var (
	ErrNoDataInPayload           = errors.New("no data in payload")
	ErrFailedToDecodePayloadData = errors.New("failed to decode payload data")
)

type Payload map[string]interface{}

func (p Payload) Version() int64 {
	return VERSION
}

func (p Payload) Get(key string) (interface{}, bool) {
	if value, ok := p[key]; ok {
		return value, true
	}
	return nil, false
}

type mockParser struct{}

func NewMockParser() *mockParser {
	return &mockParser{}
}

func (p *mockParser) SupportsVersion() int64 {
	return VERSION
}

func (p *mockParser) Encode(data interface {
	GetMap() (map[string]interface{}, error)
}, sessionID string) (map[string]interface{}, error) {

	out := make(map[string]interface{})
	mapData, err := data.GetMap()
	if err != nil {
		return nil, err
	}

	mapData["sessionID"] = sessionID
	out[DATA_PAYLOAD_TAG] = mapData
	return out, nil
}

func (p *mockParser) Decode(data map[string]interface{}) (interface {
	Version() int64
	Get(string) (interface{}, bool)
}, error) {

	mapData, ok := data[DATA_PAYLOAD_TAG]
	if !ok {
		return nil, ErrNoDataInPayload
	}

	mappedData, ok := mapData.(map[string]interface{})
	if !ok {
		return nil, ErrFailedToDecodePayloadData
	}

	var payload Payload = mappedData
	return &payload, nil
}
