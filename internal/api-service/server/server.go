package server

import (
	v1 "github.com/abhirajranjan/gochat/internal/api-service/api/v1"
	"github.com/abhirajranjan/gochat/internal/api-service/config"
	"github.com/abhirajranjan/gochat/internal/api-service/controller"
	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewServer(logger logger.ILogger, cfg *config.ServerConf, handler model.IHandler) *gin.Engine {
	engine := gin.New()
	handleApi(engine, logger, cfg, handler)
	return engine
}

func handleApi(engine *gin.Engine, logger logger.ILogger, cfg *config.ServerConf, handler model.IHandler) {
	api := engine.Group("/api")

	apiVersion1 := v1.NewV1(logger, &cfg.Auth, handler)

	apiController := controller.NewApiVersionController(&controller.ApiVersionController{
		Logger: logger,
	})

	apiController.RegisterVersion(apiVersion1)
	apiController.Handle(api)
}
