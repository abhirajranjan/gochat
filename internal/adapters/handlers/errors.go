package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setInvalidToken(ctx *gin.Context, cause string) {
	clientErrMessage(ctx, http.StatusUnauthorized, "token", fmt.Errorf("%s", cause))
	ctx.Abort()
}

func setInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

func setBadRequest(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadRequest)
}

func setBadReqWithClientErr(ctx *gin.Context, err error) {
	clientErrMessage(ctx, http.StatusBadRequest, "domain", err)
	ctx.Abort()
}

func setForbidden(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusForbidden)
}

type errResp struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

func clientErrMessage(ctx *gin.Context, resp int, errtype string, err error) {
	ctx.JSON(resp, errResp{
		Type:  errtype,
		Error: err.Error(),
	})
}
