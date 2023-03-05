package grpcServer

import (
	"context"
	"net/http"

	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/internal/api-service/proto/loginService"
)

type MockGrpcClient struct {
}

func (m *MockGrpcClient) Run() {

}

func (m *MockGrpcClient) VerifyUser(c context.Context, req *loginService.LoginRequest) (*loginService.LoginResponse, error) {
	user := loginService.UserType{UserID: req.GetUsername(), UserRoles: []int64{}}
	status := loginService.ResponseStatusType{ErrCode: http.StatusOK, Err: ""}
	return &loginService.LoginResponse{User: &user, Status: &status}, nil
}

func NewMockGrpcClient() model.IGrpcServer {
	return &MockGrpcClient{}
}
