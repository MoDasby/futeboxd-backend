package team

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/pkg/errors"
	"github.com/modasby/futeboxd-api/services/football/pkg/pagination"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /teams", h.ListByName)
	r.HandleFunc("GET /teams/{ID}", h.FindByID)
}

func (h *Handler) ListByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	page, err := pagination.NewPageWithRequest(r)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	teams, err := h.repo.FindAll(name, page.Size, page.Index)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(teams); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.ParseUint(r.PathValue("teamID"), 10, 64)
	if err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr(
			"valor inválido para o id do time",
			400,
			"HANDLER:FOOTBALL:GET_TEAM:INVALID_TEAM_ID",
		))

		return
	}

	team, err := h.repo.FindOneById(int64(teamID))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(team); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}
