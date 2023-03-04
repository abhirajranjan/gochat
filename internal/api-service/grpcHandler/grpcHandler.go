package grpcHandler

import (
	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/internal/api-service/payload"
	"github.com/abhirajranjan/gochat/internal/api-service/payload/mockParser"
	"github.com/abhirajranjan/gochat/pkg/logger"
)

type grpcHandler struct {
	logger         logger.ILogger
	grpc           model.IGrpcServer
	payloadManager model.IPayLoadManager
}

func NewGrpcHandler(logger logger.ILogger, grpcServer model.IGrpcServer) model.IHandler {
	handler := &grpcHandler{
		logger:         logger,
		grpc:           grpcServer,
		payloadManager: payload.NewManager(logger, true),
	}

	handler.payloadManager.RegisterParser(mockParser.NewMockParser())
	return handler
}
