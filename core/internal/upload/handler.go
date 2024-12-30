package upload

import (
	"io"
	"net/http"

	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
)

type UploadHandler struct {
	usecase Usecase
}

func NewUploadHandler(uc Usecase) *UploadHandler {
	return &UploadHandler{usecase: uc}
}

func (h *UploadHandler) RegisterRoutes(r *http.ServeMux, injectUser middleware.AuthMiddleware) {
	r.HandleFunc("POST /uploads/profile-picture", injectUser(h.UploadProfilePicture, false))
}

func (h *UploadHandler) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	formFile, header, err := r.FormFile("profile_picture")
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	fileBytes := make([]byte, 512)
	if _, err := formFile.Read(fileBytes); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	contentType := http.DetectContentType(fileBytes)

	if _, err := formFile.Seek(0, io.SeekStart); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	input := File{
		Size:        header.Size,
		Content:     formFile,
		ContentType: contentType,
	}

	if err := h.usecase.Upload(r.Context(), input); err != nil {
		errors.HandleHttpError(w, err)
	}
}
