package server

import (
	"github.com/abhirajranjan/gochat/internal/api-service/config"
	"github.com/abhirajranjan/gochat/internal/api-service/middlewares/AuthMiddleware"
	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/internal/api-service/server/route"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewServer(logger logger.ILogger, cfg config.ServerConf, grpcHandler model.IHandler) *gin.Engine {
	engine := gin.New()
	api := engine.Group("/api/:ver")

	jwtauth, err := AuthMiddleware.NewJWTMiddleware(cfg.Auth, logger, grpcHandler)
	if err != nil {
		logger.Panic(errors.Wrap(err, "jwt.New"))
	}

	api.GET("/:channelid/messages", jwtauth.MiddlewareFunc(), route.GetMessageRouteHandler())
	api.POST("/messages", jwtauth.MiddlewareFunc(), route.PostMessageRouteHandler())
	api.GET("/refreshtoken", jwtauth.RefreshHandler)
	api.GET("/logout", jwtauth.LogoutHandler, route.Logout())
	api.POST("/login", jwtauth.LoginHandler)

	return engine
}
