package app

import (
	"gochat/internal/core/ports"

	"github.com/gorilla/mux"
)

func initUserRoute(router *mux.Router, handler ports.Handler) {
	router.HandleFunc("/user", handler.HandleGoogleAuth())

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(mux.MiddlewareFunc(handler.AuthMiddleware()))

	userRouter.HandleFunc("/user/messages", handler.GetUserMessages())
	userRouter.HandleFunc("/user", handler.DeleteUser())
}
