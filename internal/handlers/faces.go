package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
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
	filePath, err := pkg.HandleFileUpload(w, r)
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
	return
}
