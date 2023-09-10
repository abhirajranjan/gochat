package app

import (
	"gochat/internal/core/ports"

	"github.com/gorilla/mux"
)

// init /messages endpoint
func initMessageRoute(router *mux.Router, handler ports.Handler) {
	channelrouter := router.PathPrefix("/channel").Subrouter()
	channelrouter.Use(handler.AuthMiddleware())

	channelrouter.Handle("/{channelid}/message", handler.PostMessageInChannel()).Methods("POST")
	channelrouter.Handle("/{channelid}/join", handler.JoinChannel()).Methods("POST")
	channelrouter.Handle("/{channelid}", handler.GetMessagesFromChannel()).Methods("GET")
	channelrouter.Handle("/{channelid}", handler.DeleteChannel()).Methods("DELETE")
	channelrouter.Handle("", handler.NewChannel()).Methods("POST")
}
