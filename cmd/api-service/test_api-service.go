package main

import (
	"os"

	"github.com/abhirajranjan/gochat/internal/api-service/config"
	"github.com/abhirajranjan/gochat/internal/api-service/grpcHandler"
	"github.com/abhirajranjan/gochat/internal/api-service/grpcServer"
	"github.com/abhirajranjan/gochat/internal/api-service/server"
	"github.com/abhirajranjan/gochat/pkg/configManager"
	"github.com/abhirajranjan/gochat/pkg/configManager/yamlParser"
	"github.com/abhirajranjan/gochat/pkg/logger"
)

func main() {
	// parse configs
	yamlParser := yamlParser.NewYamlParser("config.yaml")
	configManager := configManager.NewConfigManager[config.Config]()
	configManager.RegisterConfigParser(yamlParser)
	cfg := configManager.LoadConfig()

	// logger configs
	logger := logger.NewLogger(cfg.Logger)
	// add logger writer to which logger write into
	// can be anything that implements io.Writer interface
	logger.AddWriter(os.Stdout)
	logger.InitLogger()

	// grpcserver := grpcServer.NewGrpcServer(cfg.Grpc, logger)
	// use mocking grpc if we would not like to interact with main grpc server
	grpcserver := grpcServer.NewMockGrpcClient()
	grpcserver.Run()

	// add handler to grpc which converts protos to models and vica versa
	grpchandler := grpcHandler.NewGrpcHandler(logger, grpcserver)
	srv := server.NewServer(logger, &cfg.Server, grpchandler)
	srv.Run()
}
