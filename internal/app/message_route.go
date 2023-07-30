package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

// init /messages endpoint
func initMessageRoute(route *gin.RouterGroup, handler ports.Handler) {
	// get messages of channel
	route.GET("/:channelid", handler.GetMessagesFromChannel)
	// post a new message to channel
	route.POST("/:channelid", handler.PostMessageInChannel)
}
