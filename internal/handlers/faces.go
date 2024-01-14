package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kagami/go-face"
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
	faces := r.Context().Value("faces").([]string)
  label := r.Context().Value("label").(string)
  face_id := r.Context().Value("face_id").(string)
  fmt.Printf("%s %s", label, face_id)
	_ = faces
	w.WriteHeader(http.StatusOK)
	return
}
