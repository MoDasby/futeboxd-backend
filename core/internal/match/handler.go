package match

import (
	"net/http"

	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type MatchHandler struct {
	usecase MatchUsecases
}

func NewMatchHandler(usecase MatchUsecases) *MatchHandler {
	return &MatchHandler{usecase: usecase}
}

func (h *MatchHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /popular/matches", h.listPopularMatches)
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
