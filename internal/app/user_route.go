package app

import (
	"gochat/internal/core/ports"

	"github.com/gorilla/mux"
)

func initUserRoute(router *mux.Router, handler ports.Handler) {
	router.Handle("/user", handler.HandleGoogleAuth()).Methods("POST")

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(handler.AuthMiddleware())

	userRouter.Handle("/messages", handler.GetUserMessages()).Methods("GET")
	userRouter.Handle("", handler.DeleteUser()).Methods("DELETE")
}
