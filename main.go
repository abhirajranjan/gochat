package main

import (
	"gochat/config"
	"gochat/config/parser/yamlParser"
	"gochat/internal/adapters/handlers"
	"gochat/internal/adapters/repositories/sql"
	"gochat/internal/app"
	"gochat/internal/core/services"
	"gochat/logger"
	"log"
)

func main() {
	yamlparser := yamlParser.NewYamlParser("config.yaml")

	var cfg struct {
		Logger config.LoggerConfig
		App    config.AppConfig
		Sql    config.SqlConfig
		Jwt    config.JwtConfig
	}
	if err := yamlparser.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	applogger := logger.NewLogger(cfg.Logger)

	repo, err := sql.NewMySqlRepository(cfg.Sql)
	applogger.Panic(err)

	service := services.NewService(repo)
	handler := handlers.NewHandler(cfg.Jwt, service, applogger)
	server := app.NewServer(cfg.App, handler, applogger)

	server.Start()
	//TODO: add triggers and for select for graceful shutdown of server
}
