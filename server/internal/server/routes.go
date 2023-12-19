package server

import (
	"log/slog"
	"net/http"
	"os"
	"visio/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(jsonHandler)
	loggingMiddleware := middlewares.NewLoggingMiddleware(logger)
	r.Use(loggingMiddleware.Logging)
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
		return
	})
	return r
}
