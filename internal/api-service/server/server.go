package server

import (
	v1 "github.com/abhirajranjan/gochat/internal/api-service/api/v1"
	"github.com/abhirajranjan/gochat/internal/api-service/config"
	"github.com/abhirajranjan/gochat/internal/api-service/controller"
	"github.com/abhirajranjan/gochat/internal/api-service/middlewares/JwtAuthMiddleware"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/gin-gonic/gin"
)

type IHandler interface {
	HandleLoginRequest(c *gin.Context) (interface {
		GetMap() (map[string]interface{}, error)
		GetErr() string
		GetErrCode() int64
	}, error)

	GeneratePayloadData(userData interface {
		GetMap() (map[string]interface{}, error)
	}) map[string]interface{}

	ExtractPayloadData(claims map[string]interface{}) interface {
		Version() int64
		Get(string) (interface{}, bool)
		GetSessionID() interface{}
	}

	VerifyUser(data interface {
		Version() int64
		Get(string) (interface{}, bool)
		GetSessionID() interface{}
	}, reqperm []string) bool

	LogoutUser(claims map[string]interface{}) int
}

func NewServer(logger logger.ILogger, cfg *config.ServerConf, jwtHandler IHandler) *gin.Engine {
	engine := gin.New()
	handleApi(engine, logger, cfg, jwtHandler)
	return engine
}

func handleApi(engine *gin.Engine, logger logger.ILogger, cfg *config.ServerConf, jwtHandler IHandler) {
	api := engine.Group("/api")

	auth, err := JwtAuthMiddleware.NewJWTMiddleware(&cfg.Auth, logger, jwtHandler)
	if err != nil {
		logger.Fatalf("jwt middleware failed with error: ", err)
	}

	apiVersion1 := v1.NewV1(logger, auth)

	apiController := controller.NewApiVersionController(&controller.ApiVersionController{
		Logger: logger,
	})

	apiController.RegisterVersion(apiVersion1)
	apiController.Handle(api)
}
