package v1

import (
	"net/http"

	"github.com/abhirajranjan/gochat/pkg/logger"
	"github.com/gin-gonic/gin"
)

type IAuth interface {
	// checks if request have valid permission specified in permissions
	CheckIfValidAuth(permissions []string) gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
	Logout() gin.HandlerFunc
	Login() gin.HandlerFunc
}

const VERSION = "v1"

type v1 struct {
	// hold version v1 api group
	group *gin.RouterGroup

	// logger instance
	logger logger.ILogger

	// authentication instance
	Auth IAuth
}

// V1 factory method to generate new instance of api version v1
//
// if calling independently (without controller), explictly call Handle function
// to register v1 group to gin group
func NewV1(logger logger.ILogger, Auth IAuth) *v1 {
	return &v1{logger: logger, Auth: Auth}
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
	v.group.GET("/:channelid/messages", v.Auth.CheckIfValidAuth([]string{"user"}), GetMessageRouteHandler())
	v.group.POST("/messages", v.Auth.CheckIfValidAuth([]string{"user"}), PostMessageRouteHandler())
	v.group.GET("/refreshtoken", v.Auth.RefreshToken())
	v.group.GET("/logout", v.Auth.Logout(), Logout())
	v.group.POST("/login", v.Auth.Login())
}

func home() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
	}
}
