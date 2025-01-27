package handler

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/internal/user/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type Handler struct {
	usecases user.Usecases
}

func NewUserHandler(usecases user.Usecases) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux, injectUser middleware.AuthMiddleware) {
	r.HandleFunc("GET /users", injectUser(h.getCurrentUser, false))
	r.HandleFunc("POST /users", h.createUser)
	r.HandleFunc("PATCH /users", injectUser(h.editUser, false))
	r.HandleFunc("PATCH /users/password", injectUser(h.changePassword, false))
	r.HandleFunc("POST /users/recover", h.recoverPassword)
	r.HandleFunc("PATCH /users/recover", h.resetPassword)
}

func (h *Handler) getCurrentUser(w http.ResponseWriter, r *http.Request) {

	user, err := h.usecases.GetCurrentUser(r.Context())
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, user); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr("corpo de requisição inválido", 400, "HANDLER:USER:CREATE:INVALID_BODY"))

		return
	}

	if err := h.usecases.CreateUser(r.Context(), &input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *Handler) editUser(w http.ResponseWriter, r *http.Request) {
	var input dto.EditUser

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr("corpo de requisição inválido", 400, "HANDLER:USER:EDIT:INVALID_BODY"))

		return
	}

	if err := h.usecases.EditUser(r.Context(), &input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *Handler) changePassword(w http.ResponseWriter, r *http.Request) {
	var input dto.ChangePassword

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		err = errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:UPDATE_PASSWORD:INVALID_BODY",
		)

		errors.HandleHttpError(w, err)

		return
	}

	if err := h.usecases.ChangePassword(r.Context(), &input); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *Handler) recoverPassword(w http.ResponseWriter, r *http.Request) {
	var body dto.RecoverPassword

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err = errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:USER:RECOVER:INVALID_BODY",
		)

		errors.HandleHttpError(w, err)
		return
	}

	if err := h.usecases.RecoverPassword(r.Context(), &body); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var body dto.ResetPassword

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err = errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:USER:RESET_PASSWORD:INVALID_BODY",
		)

		errors.HandleHttpError(w, err)
		return
	}

	if err := h.usecases.ResetPassword(r.Context(), &body); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}
