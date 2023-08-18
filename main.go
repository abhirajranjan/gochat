package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gochat/config"
	"gochat/config/parser"
	"gochat/internal/adapters/handlers"
	"gochat/internal/adapters/repositories/sqlite"
	"gochat/internal/app"
	"gochat/internal/core/services"
	"gochat/logger"
)

func main() {
	var cfg struct {
		Logger config.LoggerConfig
		App    config.AppConfig
		Sql    config.SqlConfig
		Jwt    config.JwtConfig
	}
	if err := parser.Load(&cfg, "config.yaml"); err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded config: %+v\n", cfg)

	applogger := logger.NewLogger(cfg.Logger)
	applogger.AddWriter(os.Stdout)
	applogger.InitLogger()

	repo, err := sqlite.NewSqliteRepository(cfg.Sql)
	if err != nil {
		applogger.Panic(err)
	}

	service := services.NewService(repo)
	handler := handlers.NewHandler(cfg.Jwt, service, applogger)
	server := app.NewServer(cfg.App, handler)

	server.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	applogger.Infof("stopping server: %s", sig)
	server.Stop(context.Background())
}
