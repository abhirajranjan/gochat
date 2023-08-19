package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func initUserRoute(router *gin.Engine, handler ports.Handler) {
	router.GET("/user/messages", handler.AuthMiddleware, handler.GetUserMessages)
	router.DELETE("/user", handler.AuthMiddleware, handler.DeleteUser)
	router.POST("/user", handler.HandleGoogleAuth)
}
