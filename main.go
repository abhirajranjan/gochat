package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gochat/config"
	"gochat/config/parser"
	"gochat/internal/adapters/handlers"
	"gochat/internal/adapters/repositories/postgres"
	"gochat/internal/app"
	"gochat/internal/core/services"
	"gochat/logger"
)

type centralcfg struct {
	Logger config.LoggerConfig
	App    config.AppConfig
	Sql    config.SqlConfig
	Jwt    config.JwtConfig
}

func main() {
	cfg := defaultCfg()
	logger.AddTextLogger()

	if err := parser.Load(&cfg); err != nil {
		slog.Warn("parser.Load", "error", err)
	}

	slog.Debug("config loaded", slog.Any("cfg", cfg))

	repo, err := postgres.NewPostgresRepository(cfg.Sql)
	if err != nil {
		slog.Error("postgres.NewPostgresRepository", "error", err)
		return
	}

	service := services.NewService(repo)
	handler := handlers.NewHandler(cfg.Jwt, service)
	server := app.NewServer(cfg.App, handler)

	app.Start(server)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	slog.Info("stopping server: %s", sig)
	app.Stop(server, context.Background())
}

func defaultCfg() centralcfg {
	return centralcfg{
		Logger: config.LoggerConfig{
			Name:    "gochat-dev",
			Level:   "debug",
			Isdev:   true,
			Encoder: "console",
		},
		App: config.AppConfig{
			Addr: "localhost",
			Port: "80",
		},
		Sql: config.SqlConfig{
			SqlTimeout: 5 * time.Second,
		},
		Jwt: config.JwtConfig{
			Key:    "test",
			Expiry: 4 * time.Hour,
		},
	}
}
