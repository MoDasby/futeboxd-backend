package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/internal/session/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/cookies"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type sessionHandler struct {
	usecase session.SessionUsecases
	cfg     *config.Cookies
}

func NewsessionHandler(usecase session.SessionUsecases, cfg *config.Cookies) *sessionHandler {
	return &sessionHandler{
		usecase: usecase,
		cfg:     cfg,
	}
}

func (h *sessionHandler) RegisterRoutes(r *http.ServeMux, injectUser middleware.AuthMiddleware) {
	r.HandleFunc("POST /sessions", h.login)
	r.HandleFunc("DELETE /sessions", injectUser(h.logout, false))
}

func (h *sessionHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginInput dto.Login

	if err := json.NewDecoder(r.Body).Decode(&loginInput); err != nil {
		err = &errors.HTTPErr{
			Msg:        "Corpo de requisição inválido",
			Code:       http.StatusBadRequest,
			Context:    "SESSION:HANDLER:LOGIN:INVALID_BODY",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  r.Context().Value("traceID").(string),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}

		errors.HandleHttpError(r.Context(), w, err)
	}

	output, err := h.usecase.Login(r.Context(), &loginInput)
	if err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	cookie := cookies.CreateSessionCookie(h.cfg, output.Token, output.ExpiresAt)

	http.SetCookie(w, cookie)

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}
}

func (h *sessionHandler) logout(w http.ResponseWriter, r *http.Request) {

	if err := h.usecase.Logout(r.Context()); err != nil {
		errors.HandleHttpError(r.Context(), w, err)

		return
	}

	cookie := cookies.DeleteSessionCookie(h.cfg)

	http.SetCookie(w, cookie)

	w.WriteHeader(201)
}
