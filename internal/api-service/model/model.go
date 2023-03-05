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
	HandleLoginRequest(*gin.Context) (ILoginResponse, error)
	GeneratePayloadData(ILoginResponse) map[string]interface{}
	ExtractPayloadData(map[string]interface{}) IPayloadData
	VerifyUser(IPayloadData) bool
	LogoutUser(map[string]interface{}) int
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type IGrpcServer interface {
	Run()
	VerifyUser(context.Context, *loginService.LoginRequest) (*loginService.LoginResponse, error)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type IPayloadData interface {
	Version() int64
}

type IParser interface {
	SupportsVersion() int64
	// if inplace is true then add additional fields to input map else generate output map
	Encode(map[string]interface{}) (map[string]interface{}, error)
	Decode(map[string]interface{}) (IPayloadData, error)
	VerifyUser(IPayloadData) bool
	LogoutUser(map[string]interface{}) error
}

type IPayLoadManager interface {
	RegisterParser(IParser) error
	Encode(map[string]interface{}, int64) (map[string]interface{}, error)
	AddPayload(map[string]interface{}) (map[string]interface{}, error)
	Decode(map[string]interface{}, int64) (IPayloadData, error)
	VerifyUser(IPayloadData) bool
	LogoutUser(map[string]interface{}, int64) bool
	SetMinimumVersion(minimiumVersion int64) error
	GetMinimumVersion() int64
}
