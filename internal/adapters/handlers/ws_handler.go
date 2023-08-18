package handlers

import "github.com/gin-gonic/gin"

func (h *handler) HandleWS(ctx *gin.Context) {
	// upgrade connection to websocket
	conn, err := h.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// in case of error, websocket upgrade automatically
		// write the error to writer with status code
		h.Errorf("ws upgrade: %w", err)
		setInternalServerError(ctx)
	}

	if err := h.service.HandleWS(conn); err != nil {
		h.Errorf("handleWS: %w", err)
		setInternalServerError(ctx)
		return
	}
}
