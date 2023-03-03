package middlewares

import "github.com/gin-gonic/gin"

func ApiVersionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
