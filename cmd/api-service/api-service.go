package main

import (
	"flag"
	"log"

	"github.com/abhirajranjan/gochat/internal/api-service/config"
	"github.com/abhirajranjan/gochat/internal/api-service/grpcHandler"
	"github.com/abhirajranjan/gochat/internal/api-service/grpcServer"
	"github.com/abhirajranjan/gochat/internal/api-service/server"
	"github.com/abhirajranjan/gochat/pkg/logger"
)

func main() {
	flag.Parse()

	cfg, err := config.LoadConf()
	if err != nil {
		log.Fatal(err)
	}
	logger := logger.NewLogger(cfg.Logger)
	grpcServer := grpcServer.NewGrpcServer(cfg.Grpc, logger)
	grpcServer.Run()
	grpchandler := grpcHandler.NewGrpcHandler(logger, grpcServer)
	srv := server.NewServer(logger, cfg.Server, grpchandler)
	srv.Run()
}
