package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func initWebsocket(route *gin.RouterGroup, handler ports.Handler) {
	route.GET("/", handler.HandleWS)
}
