package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"visio/pkg"
	"visio/repositories"

	"github.com/Kagami/go-face"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	_ = current_user
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	f, header, err := r.FormFile("face")
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	file_extension := pkg.GetFileExtention(header)
	if file_extension == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body: Malformed file name"))
		return
	}
	new_file_id := uuid.NewString()
	file_path := "faces/" + new_file_id + "." + file_extension
	dst, err := os.Create(file_path)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, f); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	recognized_faces, err := h.recognizer.RecognizeFile(file_path)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(recognized_faces) == 0 {
		w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("No face detected on image"))
		return
	}
	if len(recognized_faces) > 1 {
		w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Image must contain only one recognizable face"))
		return
	}
  recognized_face := recognized_faces[0]
  new_face_id := uuid.NewString()
  err = h.faces_repo.DeleteFace()
	w.WriteHeader(http.StatusCreated)
	return
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
