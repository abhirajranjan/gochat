package app

import (
	"gochat/internal/core/ports"

	"github.com/gorilla/mux"
)

func initUserRoute(router *mux.Router, handler ports.Handler) {
	router.Handle("/user", handler.HandleGoogleAuth())

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(handler.AuthMiddleware())

	userRouter.Handle("/messages", handler.GetUserMessages())
	userRouter.Handle("", handler.DeleteUser())
}
