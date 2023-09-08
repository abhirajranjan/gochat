package app

import (
	"context"
	"fmt"
	"gochat/config"
	"gochat/internal/core/ports"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// create a new server
func NewServer(cfg config.AppConfig, srvhandler ports.Handler) *http.Server {
	router := mux.NewRouter()
	server := http.Server{
		Addr:    cfg.Addr + ":" + cfg.Port,
		Handler: router,
	}

	initUserRoute(router, srvhandler)
	initMessageRoute(router, srvhandler)
	initwebSocket(router, srvhandler)

	return &server
}

func Start(server *http.Server) {
	go func() {
		slog.Info("server started", "addr", server.Addr)
		// always returns non nil error
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("server listen: %s", err))
		}
	}()
}

func Stop(server *http.Server, ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("server listen", err)
	}
}
