package main

import (
	"gochat/config"
	"gochat/config/parser/yamlParser"
	"gochat/internal/adapters/handlers"
	"gochat/internal/adapters/repositories/redis"
	"gochat/internal/app"
	"gochat/logger"
	"log"
)

// func main() {
// 	// parse configs
// 	cfg, err := config.LoadConfig("config.yaml")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// logger configs
// 	logger := logger.NewLogger(cfg.Logger)
// 	// add logger writer to which logger write into
// 	// can be anything that implements io.Writer interface
// 	logger.AddWriter(os.Stdout)
// 	logger.InitLogger()

// 	routerEngine := router.Route(logger, &cfg.Server)
// 	routerEngine.Run()
// }

func main() {
	yamlparser := yamlParser.NewYamlParser("config.yaml")

	var cfg struct {
		Logger config.LoggerConfig
		App    config.AppConfig
	}
	if err := yamlparser.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	applogger := logger.NewLogger(cfg.Logger)

	redisrepo := redis.NewRedisRepository(applogger)
	handler := handlers.NewHandler(redisrepo, applogger)
	server := app.NewServer(cfg.App, handler, applogger)

	server.Start()
	//TODO: add triggers and for select for graceful shutdown of server
}
