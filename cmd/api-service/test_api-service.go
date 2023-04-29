package main

import (
	"os"

	"github.com/abhirajranjan/gochat/internal/api-service/config"
	"github.com/abhirajranjan/gochat/internal/api-service/dbBridge"
	"github.com/abhirajranjan/gochat/internal/api-service/jwtHandler"
	"github.com/abhirajranjan/gochat/internal/api-service/mockGrpcServer"
	"github.com/abhirajranjan/gochat/internal/api-service/payloadManager"
	"github.com/abhirajranjan/gochat/internal/api-service/payloadManager/payloadParserV1"
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
	grpcserver := mockGrpcServer.NewMockGrpcClient()
	grpcserver.Run()

	db := dbBridge.NewDbBridge(logger, grpcserver)

	payload := payloadManager.NewManager(logger)
	parserv1 := payloadParserV1.NewV1Parser()
	payload.RegisterParser(parserv1)
	payload.SetMinimumVersion(parserv1.SupportsVersion())

	jwthandler := jwtHandler.NewJwtHandler(logger, db, payload)
	srv := server.NewServer(logger, &cfg.Server, jwthandler)
	srv.Run()
}
