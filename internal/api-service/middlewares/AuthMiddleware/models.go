package AuthMiddleware

import "github.com/gin-gonic/gin"

type IHandler interface {
	HandleLoginRequest(interface{}) (interface{}, error)
	GenerateLoginRequest(c *gin.Context) (interface{}, error)
	GeneratePayloadData(ILoginResponse) map[string]interface{}
	ExtractPayloadData(map[string]interface{}) IPayloadData
}

type ILoginResponse interface {
	GetUserID() string
	GetUserRoles() []int64
	GetErr() error
	GetErrCode() int64
}

type IPayloadData interface {
	GetUserID() string
}
