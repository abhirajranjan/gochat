package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func initOAuth(router gin.IRoutes, handler ports.Handler) {
	router.Use(handler.HandleGoogleAuth)
}
