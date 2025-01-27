package handler

import (
	"net/http"

	"github.com/modasby/futeboxd-backend/core/internal/match"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type MatchHandler struct {
	usecase match.MatchUsecases
}

func NewMatchHandler(usecase match.MatchUsecases) *MatchHandler {
	return &MatchHandler{usecase: usecase}
}

func (h *MatchHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /popular/matches", h.listPopularMatches)
	r.HandleFunc("GET /popular/matches/{id}", func(w http.ResponseWriter, r *http.Request) {
		matchID, err := utils.ParseIntValue(r.PathValue("id"))
		if err != nil {
			errors.HandleHttpError(w, err)

			return
		}

		output, err := h.usecase.GetMatchStats(r.Context(), matchID)
		if err != nil {
			errors.HandleHttpError(w, err)

			return
		}

		if err := utils.SendJSON(w, output); err != nil {
			errors.HandleHttpError(w, err)

			return
		}
	})
}

func (h *MatchHandler) listPopularMatches(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.WithRequest(r)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	output, err := h.usecase.ListPopularMatches(r.Context(), page)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	if err := utils.SendJSON(w, output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}
