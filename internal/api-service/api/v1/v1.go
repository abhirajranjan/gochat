package v1

import (
	"net/http"

	"github.com/abhirajranjan/gochat/internal/api-service/middlewares/AuthMiddleware"
	"github.com/abhirajranjan/gochat/internal/api-service/model"
	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const VERSION = "v1"

type v1 struct {
	// hold version v1 api group
	group *gin.RouterGroup

	// logger instance
	logger logger.ILogger

	// authentication cfg
	authcfg *AuthMiddleware.AuthConf

	// handler for authentication to help in navigation
	handler model.IHandler
}

// V1 factory method to generate new instance of api version v1
//
// if calling independently (without controller), explictly call Handle function
// to register v1 group to gin group
func NewV1(logger logger.ILogger, authcfg *AuthMiddleware.AuthConf, handler model.IHandler) *v1 {
	return &v1{logger: logger, authcfg: authcfg, handler: handler}
}

// returns the supported version
func (v *v1) GetVersion() string {
	return VERSION
}

// create the group for the version from another group
//
// handler functions can be added as a middlewares in the group
func (v *v1) Handle(group *gin.RouterGroup, handler ...gin.HandlerFunc) {
	v.group = group.Group(VERSION)
	v.group.Use(handler...)
	v.group.GET("/", home())
	v.route()
}

// add routes to v.group
func (v *v1) route() {
	jwtauth, err := AuthMiddleware.NewJWTMiddleware(v.authcfg, v.logger, v.handler)
	if err != nil {
		v.logger.Panic(errors.Wrap(err, "jwt.New"))
	}

	v.group.GET("/:channelid/messages", jwtauth.MiddlewareFunc(), GetMessageRouteHandler())
	v.group.POST("/messages", jwtauth.MiddlewareFunc(), PostMessageRouteHandler())
	v.group.GET("/refreshtoken", jwtauth.RefreshHandler)
	v.group.GET("/logout", jwtauth.LogoutHandler, Logout())
	v.group.POST("/login", jwtauth.LoginHandler)
}

func home() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
	}
}
