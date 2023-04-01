package dbHandler

import (
	"context"

	proto "github.com/abhirajranjan/gochat/internal/api-service/proto/loginService"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

type IGrpcServer interface {
	GetUser(context.Context, *proto.LoginRequest) (*proto.LoginResponse, error)
	GetUserbyUserID(context.Context, string) (*proto.LoginResponse, error)
}

type session struct {
	Status int
	UserID string
}

type dbHandler struct {
	logger   logger.ILogger
	grpcConn IGrpcServer
	session  map[string]session
}

func NewDbHandler(logger logger.ILogger, grpcserver IGrpcServer) *dbHandler {
	return &dbHandler{logger: logger, grpcConn: grpcserver, session: make(map[string]session)}
}

func (db *dbHandler) GetUser(username string, password string) (interface {
	GetMap() (map[string]interface{}, error)
	GetErr() string
	GetErrCode() int64
}, error) {

	protoreq := proto.LoginRequest{Username: username, Password: password}
	protores, err := db.grpcConn.GetUser(context.Background(), &protoreq)

	if err != nil {
		db.logger.Error("grpc connection failed to establish: %s", err)
		return nil, ErrInternalServer
	}

	userData := protoToUser(protores)
	return userData, nil
}

func (db *dbHandler) GetUserRolesFromSession(sessionID string) (interface{ Has(string) bool }, bool) {
	session, ok := db.session[sessionID]
	if !ok {
		return nil, false
	}

	protores, err := db.grpcConn.GetUserbyUserID(context.Background(), session.UserID)
	if err != nil {
		db.logger.Error("grpc connection failed to establish: %s", err)
		return nil, false
	}

	userData := protoToUser(protores)
	return userData.User.UserRoles, true
}

func (db *dbHandler) Logout(sessionID string) error {
	delete(db.session, sessionID)
	return nil
}

func (db *dbHandler) GenerateSessionID(userdata interface{}) string {
	user, ok := userdata.(*userData)
	if !ok {
		return ""
	}

	id := uuid.New().String()
	db.session[id] = session{
		UserID: user.User.UserID,
		Status: 0,
	}

	return id
}

func (db *dbHandler) ActivateSessionByID(id string) bool {
	f := db.session[id]
	f.Status = 1
	return true
}

func protoToUser(lr *proto.LoginResponse) *userData {
	return &userData{
		User: user{
			UserID:    lr.User.UserID,
			UserRoles: lr.User.UserRoles,
		},
		Status: status{
			Err:     lr.Status.Err,
			ErrCode: lr.Status.ErrCode,
		},
	}
}
