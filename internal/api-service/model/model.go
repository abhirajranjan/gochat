package model

import (
	"context"

	"github.com/abhirajranjan/gochat/internal/api-service/proto/loginService"
	"github.com/gin-gonic/gin"
)

type ILoginRequest interface {
	GetUsername() string
	SetUsername(string)

	GetPassword() string
	SetPassword(string)
}

type ILoginResponse interface {
	GetUserID() string
	GetUserRoles() []int64
	GetErr() error
	GetErrCode() int64
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type IHandler interface {
	HandleLoginRequest(ILoginRequest) (ILoginResponse, error)
	GenerateLoginRequest(c *gin.Context) (ILoginRequest, error)
	GeneratePayloadData(ILoginResponse) map[string]interface{}
	ExtractPayloadData(claims map[string]interface{}) IPayloadData
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type IGrpcServer interface {
	Run()
	VerifyUser(ctx context.Context, loginreq *loginService.LoginRequest) (res *loginService.LoginResponse, err error)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type IPayloadData interface {
	Version() int64
}

type IParser interface {
	SupportsVersion() int64
	Encode(map[string]interface{}, bool) (map[string]interface{}, error)
	Decode(map[string]interface{}) (IPayloadData, error)
	To_Proto(IPayloadData) interface{}
}

type IPayLoadManager interface {
	RegisterParser(parser IParser) error
	Encode(data map[string]interface{}, version int64) (map[string]interface{}, error)
	AddPayload(data map[string]interface{}) (map[string]interface{}, error)
	Decode(data map[string]interface{}, version int64) (IPayloadData, error)
	To_Proto(data IPayloadData) interface{}
}
