package app

import (
	"context"
	"fmt"
	"gochat/config"
	"gochat/internal/core/ports"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ginServer struct {
	handler    ports.Handler
	engine     *gin.Engine
	httpServer *http.Server
}

func NewServer(config config.AppConfig, srvhandler ports.Handler) *ginServer {
	ginRouter := gin.New()
	server := ginServer{
		handler: srvhandler,
		engine:  ginRouter,

		httpServer: &http.Server{
			Addr:    config.Addr + ":" + config.Port,
			Handler: ginRouter,
		},
	}

	initUserRoute(ginRouter, server.handler)
	initMessageRoute(ginRouter, server.handler)

	return &server
}

func (s *ginServer) Start() {
	go func() {
		// always returns non nil error
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("server listen: %w", err))
		}
	}()
}

func (s *ginServer) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Println("server listen", err)
	}
}
