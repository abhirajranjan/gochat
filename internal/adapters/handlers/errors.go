package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUnauthorised(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

func setInvalidToken(ctx *gin.Context) {
	ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s", "invalid token"))
}

func setInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func setBadRequest(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadRequest)
}

func setBadRequestWithErr(ctx *gin.Context, err error) {
	ctx.AbortWithError(http.StatusBadRequest, err)
}

func setForbidden(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusForbidden)
}
