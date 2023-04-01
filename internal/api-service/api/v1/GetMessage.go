package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMessageRouteHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "message")
	}
}
