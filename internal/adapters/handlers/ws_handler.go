package handlers

import (
	"log/slog"
	"net/http"
)

func (h *handler) HandleWS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := h.wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("ws upgrade: %w", err)
			http.Error(w, "", http.StatusInternalServerError)
		}

		if err := h.service.HandleWS(conn); err != nil {
			slog.Error("handleWS: %w", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
