package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func initWebsocket(route gin.IRoutes, handler ports.Handler) {
	route.Any("/", handler.HandleWS)
}
