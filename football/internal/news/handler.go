package news

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-api/services/football/pkg/errors"
	"github.com/modasby/futeboxd-api/services/football/pkg/pagination"
)

type NewsHandler struct {
	repo NewsRepo
}

func NewNewsHandler(repo NewsRepo) *NewsHandler {
	return &NewsHandler{repo: repo}
}

func (h *NewsHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		page, err := pagination.NewPageWithRequest(r)
		if err != nil {
			errors.HandleHttpError(w, err)
			return
		}

		output, err := h.repo.GetRecentNews(r.Context(), page.Size, page.Index)
		if err != nil {
			errors.HandleHttpError(w, err)

			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(output); err != nil {
			errors.HandleHttpError(w, err)

			return
		}
	})
}
