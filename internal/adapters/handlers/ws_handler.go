package handlers

import "net/http"

func (h *handler) HandleWS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// upgrade connection to websocket
		conn, err := h.wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			// in case of error, websocket upgrade automatically
			// write the error to writer with status code
			h.Errorf("ws upgrade: %w", err)
			http.Error(w, "", http.StatusInternalServerError)
		}

		if err := h.service.HandleWS(conn); err != nil {
			h.Errorf("handleWS: %w", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
