package app

import (
	"gochat/internal/core/ports"

	"github.com/gorilla/mux"
)

func initwebSocket(router *mux.Router, handler ports.Handler) {
	wsRouter := router.PathPrefix("/ws").Subrouter()
	wsRouter.HandleFunc("", handler.HandleWS())
}
