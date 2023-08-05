package app

import (
	"context"
	"gochat/config"
	"gochat/internal/core/ports"
	"gochat/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ginServer struct {
	logger     logger.ILogger
	handler    ports.Handler
	engine     *gin.Engine
	httpServer *http.Server
}

func NewServer(config config.AppConfig, srvhandler ports.Handler, logger logger.ILogger) *ginServer {
	ginRouter := gin.New()
	server := ginServer{
		handler: srvhandler,
		engine:  ginRouter,

		httpServer: &http.Server{
			Addr:    config.Addr,
			Handler: ginRouter,
		},
	}

	userGroup := server.engine.Group("/user", server.handler.AuthMiddleware)
	messagesGroup := server.engine.Group("/messages", server.handler.AuthMiddleware)
	oauth := server.engine.Group("/oauth")
	ws := server.engine.Group("/ws")

	initUserRoute(userGroup, server.handler)
	initMessageRoute(messagesGroup, server.handler)
	initWebsocket(ws, server.handler)
	initOAuth(oauth, server.handler)

	return &server
}

func (s *ginServer) Start() {
	go func() {
		// always returns non nil error
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Errorf("server listen: %w", err)
		}
	}()
}

func (s *ginServer) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.WarnMsg("server listen", err)
	}
}
