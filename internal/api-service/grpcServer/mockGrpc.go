package grpcServer

import (
	"context"
	"net/http"

	"github.com/abhirajranjan/gochat/internal/api-service/proto/loginService"
)

type MockGrpcClient struct {
}

func (m *MockGrpcClient) Run() {

}

func (m *MockGrpcClient) GetUser(c context.Context, req *loginService.LoginRequest) (*loginService.LoginResponse, error) {
	user := loginService.UserType{UserID: req.GetUsername(), UserRoles: []string{"user"}}
	status := loginService.ResponseStatusType{ErrCode: http.StatusOK, Err: ""}
	return &loginService.LoginResponse{User: &user, Status: &status}, nil
}

func (m *MockGrpcClient) GetUserbyUserID(c context.Context, userid string) (*loginService.LoginResponse, error) {
	user := loginService.UserType{UserID: "abhiraj", UserRoles: []string{"user"}}
	status := loginService.ResponseStatusType{ErrCode: http.StatusOK, Err: ""}
	return &loginService.LoginResponse{User: &user, Status: &status}, nil
}

func NewMockGrpcClient() *MockGrpcClient {
	return &MockGrpcClient{}
}
