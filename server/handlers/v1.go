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

func (h *FacesHandlerv1) GetFace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	_ = current_user
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error(err)
		err = pkg.RespondToBadRequest(w, "INVALID BODY")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if _, ok := r.MultipartForm.File["face"]; !ok {
		err = pkg.RespondToBadRequest(w, "'face' FIELD NOT FOUND")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
		err = pkg.RespondToBadRequest(w, "INVALID FILE NAME")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if file_extension != "jpg" {
		err = pkg.RespondToBadRequest(w, "UNSUPPORTED FORMAT")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	new_file_id := uuid.NewString()
	file_path := os.Getenv("UPLOAD_DIR") + new_file_id + "." + file_extension
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
		err = pkg.RespondToBadRequest(w, "NO FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(recognized_faces) > 1 {
		err = pkg.RespondToBadRequest(w, "MORE THAN ONE FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	descriptor, err := json.Marshal(recognized_faces[0].Descriptor)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(
		map[string]string{
			"descriptor": string(descriptor),
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
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
		err = pkg.RespondToBadRequest(w, "INVALID BODY")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
		err = pkg.RespondToBadRequest(w, "INVALID FILE NAME")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if file_extension != "jpg" {
		err = pkg.RespondToBadRequest(w, "UNSUPPORTED FORMAT")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	new_file_id := uuid.NewString()
	file_path := os.Getenv("UPLOAD_DIR") + new_file_id + "." + file_extension
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
		err = pkg.RespondToBadRequest(w, "NO FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(recognized_faces) > 1 {
		err = pkg.RespondToBadRequest(w, "MORE THAN ONE FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
	if _, err := uuid.Parse(face_id); err != nil {
		err = pkg.RespondToBadRequest(w, "'face' PARAMETER IS NOT A UUID")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	targetted_face, err := h.faces_repo.GetFaceById(face_id, current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			err = pkg.RespondToBadRequest(w, "FACE NOT FOUND")
			if err != nil {
				h.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error(err)
		err = pkg.RespondToBadRequest(w, "INVALID BODY")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	f, header, err := r.FormFile("face")
	if err != nil {
		h.logger.Error(err)
		err = pkg.RespondToBadRequest(w, "'face' FIELD NOT FOUND")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	defer f.Close()
	file_extension := pkg.GetFileExtention(header)
	if file_extension == "" {
		err = pkg.RespondToBadRequest(w, "INVALID FILE NAME")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if file_extension != "jpg" {
		err = pkg.RespondToBadRequest(w, "UNSUPPORTED FORMAT")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	new_file_id := uuid.NewString()
	file_path := os.Getenv("UPLOAD_DIR") + new_file_id + "." + file_extension
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
		err = os.Remove(file_path)
		if err != nil {
			h.logger.Warn("Failed to remove a file: ", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = os.Remove(file_path)
	if err != nil {
		h.logger.Warn("Failed to remove a file: ", err.Error())
	}
	if len(recognized_faces) == 0 {
		err = pkg.RespondToBadRequest(w, "NO FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(recognized_faces) > 1 {
		err = pkg.RespondToBadRequest(w, "MORE THAN ONE FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
	if _, err := uuid.Parse(face_id); err != nil {
		err = pkg.RespondToBadRequest(w, "'face' PARAMETER IS NOT A UUID")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
	if _, err := uuid.Parse(payload.Face1); err != nil {
		err = pkg.RespondToBadRequest(w, "'face_1' PARAMETER IS NOT A UUID")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if _, err := uuid.Parse(payload.Face2); err != nil {
		err = pkg.RespondToBadRequest(w, "'face_2' PARAMETER IS NOT A UUID")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	face1, err := h.faces_repo.GetFaceById(payload.Face1, current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			err = pkg.RespondToBadRequest(w, "face_1 NOT FOUND")
			if err != nil {
				h.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	face2, err := h.faces_repo.GetFaceById(payload.Face2, current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			err = pkg.RespondToBadRequest(w, "face_2 NOT FOUND")
			if err != nil {
				h.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
	current_user, ok := r.Context().Value("current_user").(map[string]string)
	_ = current_user
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := r.ParseMultipartForm(10 << 50)
	f1, f1_header, err := r.FormFile("face1")
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f1.Close()
	file_extension := pkg.GetFileExtention(f1_header)
	if file_extension == "" {
		err = pkg.RespondToBadRequest(w, "INVALID FILE NAME ON 'face1'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if file_extension != "jpg" {
		err = pkg.RespondToBadRequest(w, "UNSUPPORTED FILE FORMAT ON 'face1'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	new_file_id := uuid.NewString()
	file_path := os.Getenv("UPLOAD_DIR") + new_file_id + "." + file_extension
	dst, err := os.Create(file_path)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, f1); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	recognized_faces, err := h.recognizer.RecognizeFile(file_path)
	if err != nil {
		h.logger.Error(err)
		err = os.Remove(file_path)
		if err != nil {
			h.logger.Warn("Failed to remove a file ", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(recognized_faces) == 0 {
		err = pkg.RespondToBadRequest(w, "NO FACE DETECTED ON 'face1'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(recognized_faces) > 1 {
		err = pkg.RespondToBadRequest(w, "MORE THAN ONE FACE DETECTED ON 'face1'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	face1_descriptor := recognized_faces[0].Descriptor
	f2, f2_header, err := r.FormFile("face2")
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f2.Close()
	file_extension = pkg.GetFileExtention(f2_header)
	if file_extension == "" {
		err = pkg.RespondToBadRequest(w, "INVALID FILE NAME on 'face2'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if file_extension != "jpg" {
		err = pkg.RespondToBadRequest(w, "UNSUPPORTED FORMAT ON 'face2'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	new_file_id = uuid.NewString()
	file_path = os.Getenv("UPLOAD_DIR") + new_file_id + "." + file_extension
	dst, err = os.Create(file_path)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, f2); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	recognized_faces, err = h.recognizer.RecognizeFile(file_path)
	if err != nil {
		h.logger.Error(err)
		err = os.Remove(file_path)
		if err != nil {
			h.logger.Warn("Failed to remove a file ", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(recognized_faces) == 0 {
		err = pkg.RespondToBadRequest(w, "NO FACE DETECTED ON 'face2'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(recognized_faces) > 1 {
		err = pkg.RespondToBadRequest(w, "MORE THAT ONE FACE DETECTED ON 'face2'")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	face2_descriptor := recognized_faces[0].Descriptor
	euclidean_distance := face.SquaredEuclideanDistance(face1_descriptor, face2_descriptor)
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

func (h *FacesHandlerv1) CompareFacesMixt(w http.ResponseWriter, r *http.Request) {
	current_user, ok := r.Context().Value("current_user").(map[string]string)
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
	face_id_value := r.FormValue("face_id")
	if face_id_value == "" {
		err = pkg.RespondToBadRequest(w, "'face_id' PARAMETER NOT FOUND")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if _, err := uuid.Parse(face_id_value); err != nil {
		err = pkg.RespondToBadRequest(w, "'face_id' NOT A UUID")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	db_face, err := h.faces_repo.GetFaceById(face_id_value, current_user["id"])
	if err != nil {
		if err == sql.ErrNoRows {
			err = pkg.RespondToBadRequest(w, "FACE NOT FOUND")
			if err != nil {
				h.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var db_face_descriptor *face.Descriptor
	err = json.Unmarshal([]byte(db_face.Descriptor), &db_face_descriptor)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f, header, err := r.FormFile("face")
	if err != nil {
		h.logger.Error(err)
		err = pkg.RespondToBadRequest(w, "'face' FIELD NOT FOUND")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	defer f.Close()
	file_extension := pkg.GetFileExtention(header)
	if file_extension == "" {
		err = pkg.RespondToBadRequest(w, "INVALID FILE NAME")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if file_extension != "jpg" {
		err = pkg.RespondToBadRequest(w, "UNSUPPORTED FORMAT")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	new_file_id := uuid.NewString()
	file_path := os.Getenv("UPLOAD_DIR") + new_file_id + "." + file_extension
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
		err = os.Remove(file_path)
		if err != nil {
			h.logger.Warn("Failed to remove a file: ", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(recognized_faces) == 0 {
		err = pkg.RespondToBadRequest(w, "NO FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if len(recognized_faces) > 1 {
		err = pkg.RespondToBadRequest(w, "MORE THAN ONE FACE DETECTED")
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	image_face_descriptor := &recognized_faces[0].Descriptor
	euclidean_distance := face.SquaredEuclideanDistance(*image_face_descriptor, *db_face_descriptor)
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

func (h *FacesHandlerv1) RegisterRoutes(r chi.Router) {
	r.Get("/faces", h.GetFaces)
	r.Post("/faces", h.CreateFace)
	r.Put("/faces/{face}", h.UpdateFace)
	r.Delete("/faces/{face}", h.DeleteFace)

	r.Post("/faces/detect", h.GetFace)
	r.Post("/faces/compare", h.CompareFaces)
	r.Post("/faces/compare-images", h.CompareFacesWithUpload)
	r.Post("/faces/compare-mixt", h.CompareFacesMixt)
}
