package mockParser

import (
	"github.com/abhirajranjan/gochat/internal/api-service/model"
)

type mockParser struct{}

type payloadData map[string]interface{}

func (p *payloadData) Version() int64 {
	return 0
}

func NewMockParser() model.IParser {
	return &mockParser{}
}

func (p *mockParser) SupportsVersion() int64 {
	return 0
}

func (p *mockParser) VerifyUser(data model.IPayloadData) bool {
	return true
}

func (p *mockParser) Encode(i map[string]interface{}, inplace bool) (map[string]interface{}, error) {
	return i, nil
}

func (p *mockParser) Decode(data map[string]interface{}) (model.IPayloadData, error) {
	a := payloadData(data)
	return &a, nil
}

func (p *mockParser) LogoutUser(map[string]interface{}) error {
	return nil
}
