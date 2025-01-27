package upload

import (
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
	formFile, _, err := r.FormFile("profile_picture")
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := h.usecase.Upload(r.Context(), formFile); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}
