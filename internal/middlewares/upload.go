package middlewares

import (
	"log/slog"
	"net/http"
)

type UploadMiddleware struct {
	logger *slog.Logger
}

func NewUploadMiddleware(logger *slog.Logger) *UploadMiddleware {
	return &UploadMiddleware{
		logger: logger,
	}
}

func (m *UploadMiddleware) HandleUploads(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
      defer r.Body.Close()
			err := r.ParseMultipartForm(20 << 20)
		},
	)
}
