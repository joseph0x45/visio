package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"visio/internal/store"
	"visio/internal/types"
	"visio/pkg"

	"github.com/Kagami/go-face"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
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
	currentUser, ok := r.Context().Value("currentUser").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	faces, ok := r.Context().Value("faces").([]string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		deleteErrors := pkg.CleanupFiles(faces)
		for _, err := range deleteErrors {
			h.logger.Debug(err.Error())
		}
	}()
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
	descriptor, err := json.Marshal(recognizedFace.Descriptor)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while marshalling descriptor: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newFace := &types.Face{
		Id:         ulid.Make().String(),
		Label:      label,
		UserId:     currentUser,
		Descriptor: string(descriptor),
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

func (h *FaceHandler) CompareSavedFaces(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	subject, ok := r.Context().Value("subject").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if subject == "" {
		http.Error(w, "Missing field 'subject' in request body", http.StatusBadRequest)
		return
	}
	object, ok := r.Context().Value("object").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if object == "" {
		http.Error(w, "Missing field 'object' in request body", http.StatusBadRequest)
		return
	}
	subFace, err := h.faces.GetById(subject, currentUser)
	if err != nil {
		if errors.Is(err, types.ErrFaceNotFound) {
			http.Error(w, fmt.Sprintf("Face %s not found", subject), http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error())
		return
	}
	objFace, err := h.faces.GetById(object, currentUser)
	if err != nil {
		if errors.Is(err, types.ErrFaceNotFound) {
			http.Error(w, fmt.Sprintf("Face %s not found", subject), http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error())
		return
	}
	var subFaceDescriptor face.Descriptor
	err = json.Unmarshal([]byte(subFace.Descriptor), &subFaceDescriptor)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while unmarshalling descriptor: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var objFaceDescriptor face.Descriptor
	err = json.Unmarshal([]byte(objFace.Descriptor), &objFaceDescriptor)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while unmarshalling descriptor: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	distance := face.SquaredEuclideanDistance(objFaceDescriptor, subFaceDescriptor)
	data, err := json.Marshal(map[string]float64{
		"distance": distance,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while marshalling response: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (h *FaceHandler) DeleteFace(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("currentUser").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	faceId := chi.URLParam(r, "id")
	if faceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.faces.Delete(faceId, currentUser)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
