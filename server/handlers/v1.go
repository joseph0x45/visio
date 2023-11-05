package handlers

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"visio/repositories"

	"github.com/Kagami/go-face"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type FacesHandlerv1 struct {
	logger     *logrus.Logger
	faces_repo *repositories.FacesRepo
	recognizer *face.Recognizer
}

func NewFacesHandlerV1(logger *logrus.Logger, faces_repo *repositories.FacesRepo, rec *face.Recognizer) *FacesHandlerv1 {
	return &FacesHandlerv1{
		logger:     logger,
		faces_repo: faces_repo,
		recognizer: rec,
	}
}

func (h *FacesHandlerv1) GetFaces(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	faces, err := h.faces_repo.SelectAllFacesCreatedByUser(current_user["id"])
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(faces)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *FacesHandlerv1) CreateFace(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20+512)
	reader, err := r.MultipartReader()
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	part, err := reader.NextPart()
	if err != nil {
		if err == io.EOF {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if part.FormName() != "face" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	buf_reader := bufio.NewReader(part)
	sniff, err := buf_reader.Peek(512)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content_type := http.DetectContentType(sniff)
	if content_type != "image/jpeg" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	f, err := os.CreateTemp("", "")
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	var max_size int64 = 10 << 20
	lmt := io.MultiReader(buf_reader, io.LimitReader(part, max_size-511))
	written, err := io.Copy(f, lmt)
	if err != nil && err != io.EOF {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if written > max_size {
		os.Remove(f.Name())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *FacesHandlerv1) UpdateFace(w http.ResponseWriter, r *http.Request) {
}
func (h *FacesHandlerv1) DeleteFace(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	face_id := chi.URLParam(r, "face")
	if face_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rows_affected, err := h.faces_repo.DeleteFace(current_user["id"], face_id)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rows_affected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
func (h *FacesHandlerv1) CompareFaces(w http.ResponseWriter, r *http.Request) {
}
func (h *FacesHandlerv1) DetectFace(w http.ResponseWriter, r *http.Request) {
}

func (h *FacesHandlerv1) RegisterRoutes(r chi.Router) {
	r.Get("/v1/faces", h.GetFaces)
	r.Post("/v1/faces", h.CreateFace)
	r.Put("/v1/faces/{face}", h.UpdateFace)
	r.Delete("/v1/faces/{face}", h.DeleteFace)

	r.Post("/v1/faces/detect", nil)
	r.Post("/v1/faces/compare", nil)
}
