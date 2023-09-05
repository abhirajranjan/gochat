package app

import (
	"gochat/internal/core/ports"

	"github.com/gorilla/mux"
)

// init /messages endpoint
func initMessageRoute(router *mux.Router, handler ports.Handler) {
	channelrouter := router.PathPrefix("/channel").Subrouter()
	channelrouter.Use(handler.AuthMiddleware())

	channelrouter.Handle("", handler.NewChannel())
	channelrouter.Handle("/{channelid}", handler.GetMessagesFromChannel())
	channelrouter.Handle("/{channelid}/message", handler.PostMessageInChannel())
	channelrouter.Handle("/{channelid}/join", handler.JoinChannel())
	channelrouter.Handle("/{channelid}", handler.DeleteChannel())
}
