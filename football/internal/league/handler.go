package league

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-api/services/football/pkg/errors"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /leagues", func(w http.ResponseWriter, r *http.Request) {
		leagues, err := h.repo.ListLeagues(r.Context())
		if err != nil {
			errors.HandleHttpError(w, err)

			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(leagues); err != nil {
			errors.HandleHttpError(w, err)

			return
		}
	})
}
