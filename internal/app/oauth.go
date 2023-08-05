package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func initOAuth(router *gin.RouterGroup, handler ports.Handler) {
	router.GET("/", handler.HandleGoogleAuth)
}
