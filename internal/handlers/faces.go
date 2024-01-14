package handlers

import (
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
	w.WriteHeader(http.StatusOK)
	return
}
