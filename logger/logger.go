package logger

import (
	"log/slog"
	"net/http"
	"os"
)

func AddTextLogger() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	slog.SetDefault(slog.New(handler))
}

func RouterPathLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("route called", "url", r.Host+"/"+r.RequestURI, "method", r.Method)
		next.ServeHTTP(w, r)
	})
}
