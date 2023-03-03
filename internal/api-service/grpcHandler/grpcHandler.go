package grpcHandler

import (
	"github.com/abhirajranjan/gochat/internal/api-service/grpcServer"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/gin-gonic/gin"
)

type IHandler interface {
	HandleLoginRequest(ILoginRequest) (ILoginResponse, error)
	GenerateLoginRequest(c *gin.Context) (ILoginRequest, error)
	ExtractPayloadData(claims map[string]interface{}) IPayloadData
}

type grpcHandler struct {
	logger logger.ILogger
	grpc   grpcServer.IGrpcServer
}

func NewGrpcHandler(logger logger.ILogger, grpcServer grpcServer.IGrpcServer) IHandler {
	return &grpcHandler{
		logger: logger,
		grpc:   grpcServer,
	}
}
