package app

import (
	"gochat/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func initUserRoute(group *gin.RouterGroup, handler ports.Handler) {
	// get recent messages of user
	group.GET("/messages", handler.GetUserMessages)
}
