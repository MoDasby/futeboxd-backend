package handler

import (
	"net/http"

	"github.com/modasby/futeboxd-backend/core/internal/profile"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type profileHandler struct {
	usecase profile.ProfileUsecases
}

func NewProfileHandler(
	usecase profile.ProfileUsecases,
) *profileHandler {
	return &profileHandler{
		usecase: usecase,
	}
}

func (h *profileHandler) RegisterRoutes(r *http.ServeMux, injectUser middleware.AuthMiddleware) {
	r.HandleFunc("GET /profiles/{username}", injectUser(h.findByUsername, true))
	r.HandleFunc("POST /profiles/{username}/follow", injectUser(h.toggleFollow, false))
	r.HandleFunc("GET /profiles", injectUser(h.searchProfiles, true))
	r.HandleFunc("GET /profiles/{username}/followers", injectUser(h.ListFollowers, true))
	r.HandleFunc("GET /profiles/{username}/following", injectUser(h.ListFollowing, true))
	r.HandleFunc("GET /popular/profiles", injectUser(h.ListPopularProfiles, true))
}

func (h *profileHandler) ListPopularProfiles(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.WithRequest(r)
	if err != nil {
		errors.HandleHttpError(w, err)
		return
	}

	output, err := h.usecase.ListPopularProfiles(r.Context(), page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *profileHandler) ListFollowing(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	page, err := pagination.WithRequest(r)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	profiles, err := h.usecase.ListFollowing(r.Context(), username, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, profiles); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *profileHandler) ListFollowers(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	page, err := pagination.WithRequest(r)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	profiles, err := h.usecase.ListFollowers(r.Context(), username, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, profiles); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *profileHandler) findByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	output, err := h.usecase.FindByUsername(r.Context(), username)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *profileHandler) toggleFollow(w http.ResponseWriter, r *http.Request) {
	followingUsername := r.PathValue("username")

	output, err := h.usecase.ToggleFollow(r.Context(), followingUsername)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.WriteHeader(201)
}

func (h *profileHandler) searchProfiles(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.WithPage(
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	username := r.URL.Query().Get("username")

	output, err := h.usecase.SearchByUsername(r.Context(), username, page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}
