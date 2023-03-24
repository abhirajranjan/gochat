package config

import (
	"github.com/abhirajranjan/gochat/internal/api-service/grpcServer"
	"github.com/abhirajranjan/gochat/internal/api-service/middlewares/JwtAuthMiddleware"
	"github.com/abhirajranjan/gochat/pkg/logger"
)

type Config struct {
	AppName string                `mapstructure:"appName"`
	Server  ServerConf            `mapstructure:"server"`
	Logger  logger.LoggerConf     `mapstructure:"logger"`
	Grpc    grpcServer.GrpcConfig `mapstructure:"grpc"`
}

type ServerConf struct {
	Auth JwtAuthMiddleware.AuthConf `mapstructure:"auth"`
}
