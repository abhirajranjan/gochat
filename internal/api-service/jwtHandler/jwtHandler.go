package jwtHandler

import (
	"github.com/abhirajranjan/gochat/pkg/logger"
)

type ILoginRequest interface {
	GetUsername() string
	SetUsername(string)

	GetPassword() string
	SetPassword(string)
}

type IPayloadData interface {
	Version() int64
	Get(string) (interface{}, bool)
}

type IDbHandler interface {
	GetUser(username string, password string) (response interface {
		GetMap() (map[string]interface{}, error)
		GetErrCode() int64
		GetErr() string
	}, err error)

	GetUserRolesFromSession(string) (interface{ Has(string) bool }, bool)

	Logout(sessionID string) error

	GenerateSessionID(interface{}) string

	ActivateSessionByID(string) bool
}

type IPayLoadManager interface {
	Encode(data interface {
		GetMap() (map[string]interface{}, error)
	}, sessionID string, out map[string]interface{}) error

	Decode(data map[string]interface{}) (interface {
		Version() int64
		Get(string) (interface{}, bool)
	}, error)
}

type jwtHandler struct {
	logger         logger.ILogger
	payloadManager IPayLoadManager
	dbHandler      IDbHandler
}

func NewJwtHandler(logger logger.ILogger, dbhandler IDbHandler, manager IPayLoadManager) *jwtHandler {
	handler := &jwtHandler{
		logger:         logger,
		payloadManager: manager,
		dbHandler:      dbhandler,
	}

	return handler
}
