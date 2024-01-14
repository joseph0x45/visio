package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

type FaceHandler struct {
	logger *slog.Logger
}

func NewFaceHandler(logger *slog.Logger) *FaceHandler {
	return &FaceHandler{
		logger: logger,
	}
}

func (h *FaceHandler) SaveFace(w http.ResponseWriter, r *http.Request) {
	faces, ok := r.Context().Value("faces").([]string)
  _ = faces
	if !ok {
		h.logger.Debug(fmt.Sprintf("Coercion failed"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
