package handler

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/internal/session/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type sessionHandler struct {
	usecase session.SessionUsecases
}

func NewsessionHandler(usecase session.SessionUsecases) *sessionHandler {
	return &sessionHandler{
		usecase: usecase,
	}
}

func (h *sessionHandler) RegisterRoutes(r *http.ServeMux, injectUser middleware.AuthMiddleware) {
	r.HandleFunc("POST /sessions", h.login)
	r.HandleFunc("DELETE /sessions", injectUser(h.logout, false))
}

func (h *sessionHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginInput dto.Login

	if err := json.NewDecoder(r.Body).Decode(&loginInput); err != nil {
		err = errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:USER:LOGIN:INVALID_BODY",
		)
	}

	output, err := h.usecase.Login(r.Context(), &loginInput)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *sessionHandler) logout(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.Logout(r.Context()); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

}
