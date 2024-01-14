package middlewares

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"visio/pkg"
)

const (
	MAX_UPLOAD_SIZE = 5 * 1024 * 1024
)

type UploadMiddleware struct {
	logger *slog.Logger
}

func NewUploadMiddleware(logger *slog.Logger) *UploadMiddleware {
	return &UploadMiddleware{
		logger: logger,
	}
}

func (m *UploadMiddleware) HandleUploads(requiredImages int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
			if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
				if err.Error() == "http: request body too large" {
					http.Error(w, "Request body too large. Maximum size allowed is 5mb", http.StatusBadRequest)
					return
				}
				m.logger.Error(fmt.Sprintf("Error while parsing multipart form: %s", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			faces := r.MultipartForm.File["faces"]
			if len(faces) != requiredImages {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			facesPaths := []string{}
			for _, fileHeader := range faces {
				file, err := fileHeader.Open()
				if err != nil {
					m.logger.Error(fmt.Sprintf("Error while getting file header: %s", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer file.Close()
				buffer := make([]byte, 512)
				_, err = file.Read(buffer)
				if err != nil {
					m.logger.Error(fmt.Sprintf("Error while reading data from file: %s", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				fileType := http.DetectContentType(buffer)
				if fileType != "image/jpeg" && fileType != "image/png" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				_, err = file.Seek(0, io.SeekStart)
				if err != nil {
					m.logger.Error(fmt.Sprintf("Error while reseting file pointer: %s", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				f, err := os.CreateTemp("", "")
				if err != nil {
					m.logger.Error(fmt.Sprintf("Error while creating temp file: %s", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, err = io.Copy(f, file)
				if err != nil {
					m.logger.Error(fmt.Sprintf("Error while copying bytes to file: %s", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if fileType == "image/png" {
					jpegFile, err := pkg.PNGToJPEG(f.Name())
					if err != nil {
						m.logger.Error(err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					err = os.Remove(f.Name())
					if err != nil {
						m.logger.Debug(fmt.Sprintf("Error while deleting PNG file: %s", err.Error()))
					}
					facesPaths = append(facesPaths, jpegFile)
					continue
				}
				facesPaths = append(facesPaths, f.Name())
			}
			fmt.Printf("%v", facesPaths)
			ctx := context.WithValue(r.Context(), "faces", facesPaths)
			ctx = context.WithValue(ctx, "label", r.FormValue("label"))
			ctx = context.WithValue(ctx, "face_id", r.FormValue("face_id"))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
