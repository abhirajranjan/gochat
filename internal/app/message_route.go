package app

import (
	"gochat/internal/core/ports"

	"github.com/gorilla/mux"
)

// init /messages endpoint
func initMessageRoute(router *mux.Router, handler ports.Handler) {
	channelrouter := router.PathPrefix("/channel").Subrouter()
	channelrouter.Use(mux.MiddlewareFunc(handler.AuthMiddleware()))

	channelrouter.HandleFunc("", handler.NewChannel())
	channelrouter.HandleFunc("/:channelid", handler.GetMessagesFromChannel())
	channelrouter.HandleFunc("/:channelid/message", handler.PostMessageInChannel())
	channelrouter.HandleFunc("/:channelid/join", handler.JoinChannel())
	channelrouter.HandleFunc("/:channelid", handler.DeleteChannel())
}
