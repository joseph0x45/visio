package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"visio/internal/types"
	"visio/pkg"
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
	filePath, err, isJPEG := pkg.HandleFileUpload(w, r)
	if err != nil {
		if errors.Is(err, types.ErrFileNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(types.ErrFileNotFoundMessage))
			return
		}
		if errors.Is(err, types.ErrUnsupportedFormat) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(types.ErrUnsupportedFormatMessage))
			return
		}
		if errors.Is(err, types.ErrBodyTooLarge) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(types.ErrBodyTooLargeMessage))
			return
		}
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(filePath)
	fmt.Println(isJPEG)
	if !isJPEG {
		jpegFilePath, err := pkg.PNGToJPEG(filePath)
		if err != nil {
			h.logger.Error(err.Error())
			err = os.Remove(filePath)
			if err != nil {
				h.logger.Error(fmt.Sprintf("Failed to delete file: %s", err.Error()))
			}
			err = os.Remove(jpegFilePath)
			if err != nil {
				h.logger.Error(fmt.Sprintf("Failed to delete file: %s", err.Error()))
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = os.Remove(filePath)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Failed to delete file: %s", err.Error()))
		}
    filePath = jpegFilePath
	}
  fmt.Println(filePath)
	w.WriteHeader(http.StatusOK)
	return
}
