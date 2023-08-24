package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

// init /messages endpoint
func initMessageRoute(router *gin.Engine, handler ports.Handler) {
	router.GET("/ws", handler.AuthMiddleware, handler.HandleWS)
	router.POST("/channel", handler.NewChannel)

	channel := router.Group("/channel", handler.AuthMiddleware)
	channel.GET("/:channelid", handler.GetMessagesFromChannel)
	channel.POST("/:channelid/message", handler.PostMessageInChannel)
	channel.POST("/:channelid/join", handler.JoinChannel)
	channel.DELETE("/:channelid", handler.DeleteChannel)
}
