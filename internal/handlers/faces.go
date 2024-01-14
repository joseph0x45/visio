package handlers

import (
	"fmt"
	"github.com/Kagami/go-face"
	"log/slog"
	"net/http"
)

type FaceHandler struct {
	logger     *slog.Logger
	recognizer *face.Recognizer
}

func NewFaceHandler(logger *slog.Logger, recognizer *face.Recognizer) *FaceHandler {
	return &FaceHandler{
		logger:     logger,
		recognizer: recognizer,
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
