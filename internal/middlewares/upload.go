package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
)

const (
	MAX_UPLOAD_SIZE = 5 * 1024 * 1024
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
			r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
			if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
				if err.Error() == "http: request body too large" {
					http.Error(w, "Request body too large. Maximum size allowed is 5mb", http.StatusBadRequest)
					return
				}
				m.logger.Error(fmt.Sprintf("Error while parsing multipart form: %s", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		},
	)
}
