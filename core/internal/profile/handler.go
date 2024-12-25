package profile

import (
	"net/http"

	"github.com/modasby/futeboxd-backend/core/pkg/auth"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type profileHandler struct {
	usecase ProfileUsecases
}

func NewProfileHandler(
	usecase ProfileUsecases,
) *profileHandler {
	return &profileHandler{
		usecase: usecase,
	}
}

func (h *profileHandler) RegisterRoutes(r *http.ServeMux, injectUser auth.Middleware) {
	r.HandleFunc("GET /profiles/{username}", injectUser(h.findByUsername, true))
	r.HandleFunc("POST /profiles/{username}/follow", injectUser(h.toggleFollow, false))
	r.HandleFunc("GET /profiles", injectUser(h.searchProfiles, true))
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
