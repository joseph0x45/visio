package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
	"visio/models"
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
		return
	}
	if file_extension != "jpg" {
		w.WriteHeader(http.StatusBadRequest)
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
	err = os.Remove(file_path)
	if err != nil {
		h.logger.Warn("Failed to remove a file: ", err.Error())
	}
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
	descriptor, err := json.Marshal(recognized_faces[0].Descriptor)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	new_face_id := uuid.NewString()
	new_face := models.Face{
		Id:          new_face_id,
		CreatedBy:   current_user["id"],
		Descriptor:  string(descriptor),
		CreatedAt:   time.Now().String(),
		LastUpdated: time.Now().String(),
	}
	err = h.faces_repo.InsertFace(&new_face)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(
		map[string]string{
			"id": new_face_id,
		},
	)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}

func (h *FacesHandlerv1) UpdateFace(w http.ResponseWriter, r *http.Request) {
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
	if _, err := uuid.Parse(face_id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	targetted_face, err := h.faces_repo.GetFaceById(face_id, current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = r.ParseMultipartForm(10 << 20)
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
		return
	}
	if file_extension != "jpg" {
		w.WriteHeader(http.StatusBadRequest)
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
	err = os.Remove(file_path)
	if err != nil {
		h.logger.Warn("Failed to remove a file: ", err.Error())
	}
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(recognized_faces) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(recognized_faces) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	descriptor, err := json.Marshal(recognized_faces[0].Descriptor)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.faces_repo.UpdateFace(targetted_face.Id, current_user["id"], string(descriptor), time.Now().String())
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
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
	err := h.faces_repo.DeleteFace(face_id, current_user["id"])
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *FacesHandlerv1) CompareFaces(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var payload = new(
		struct {
			Face1 string `json:"face_1"`
			Face2 string `json:"face_2"`
		},
	)
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if payload.Face1 == "" || payload.Face2 == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	face1, err := h.faces_repo.GetFaceById(payload.Face1, current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	face2, err := h.faces_repo.GetFaceById(payload.Face2, current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var face1_descriptor *face.Descriptor
	err = json.Unmarshal([]byte(face1.Descriptor), &face1_descriptor)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var face2_descriptor *face.Descriptor
	err = json.Unmarshal([]byte(face2.Descriptor), &face2_descriptor)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	euclidean_distance := face.SquaredEuclideanDistance(*face1_descriptor, *face2_descriptor)
	data, err := json.Marshal(
		map[string]interface{}{
			"distance": euclidean_distance,
		},
	)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (h *FacesHandlerv1) CompareFacesWithUpload(w http.ResponseWriter, r *http.Request) {
}

func (h *FacesHandlerv1) CompareFacesMixt(w http.ResponseWriter, r *http.Request) {
}

func (h *FacesHandlerv1) RegisterRoutes(r chi.Router) {
	r.Get("/v1/faces", h.GetFaces)
	r.Post("/v1/faces", h.CreateFace)
	r.Put("/v1/faces/{face}", h.UpdateFace)
	r.Delete("/v1/faces/{face}", h.DeleteFace)

	r.Post("/v1/faces/compare", h.CompareFaces)
	r.Post("/v1/faces/compare-images", h.CompareFacesWithUpload)
	r.Post("/v1/faces/compare-mixt", h.CompareFacesMixt)
}
