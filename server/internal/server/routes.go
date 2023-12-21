package server

import (
	"fmt"
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
	r.Use(loggingMiddleware.SpamFilter)
	r.Use(loggingMiddleware.RequestLogger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from health")
		w.WriteHeader(http.StatusOK)
		return
	})
	return r
}
