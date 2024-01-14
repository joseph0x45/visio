package handlers

import (
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"net/http"
	"visio/internal/store"
	"visio/internal/types"
)

type FaceHandler struct {
	logger     *slog.Logger
	recognizer *face.Recognizer
	faces      *store.Faces
}

func NewFaceHandler(logger *slog.Logger, recognizer *face.Recognizer, faces *store.Faces) *FaceHandler {
	return &FaceHandler{
		logger:     logger,
		recognizer: recognizer,
		faces:      faces,
	}
}

func (h *FaceHandler) SaveFace(w http.ResponseWriter, r *http.Request) {
	faces, ok := r.Context().Value("faces").([]string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	facePath := faces[0]
	label, ok := r.Context().Value("label").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	count, err := h.faces.CountByLabel(label)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}
	recognizedFaces, err := h.recognizer.RecognizeFile(facePath)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while recognizing file: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(recognizedFaces) > 1 {
		http.Error(w, "More than one face detected on image", http.StatusBadRequest)
		return
	}
	if len(recognizedFaces) == 0 {
		http.Error(w, "No face detected on image", http.StatusBadRequest)
		return
	}
	recognizedFace := recognizedFaces[0]
	newFace := &types.Face{
		Id:         ulid.Make().String(),
		Label:      label,
		UserId:     "01HM54P0CZY4PPK9C8PGM2Z60J",
		Descriptor: fmt.Sprintf("%v", recognizedFace.Descriptor),
	}
	err = h.faces.Save(newFace)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
