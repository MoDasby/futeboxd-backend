package match

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/pkg/errors"
	"github.com/modasby/futeboxd-api/services/football/pkg/pagination"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /matches/{matchID}", h.GetMatch)
	r.HandleFunc("GET /matches", h.ListMatches)
	r.HandleFunc("GET /matches/{matchID}/summary", h.GetMatchSummary)
	r.HandleFunc("POST /matches/batch", h.FindMatchesBatch)
}

func (h *Handler) GetMatchSummary(w http.ResponseWriter, r *http.Request) {
	matchID, err := strconv.ParseUint(r.PathValue("matchID"), 10, 64)
	if err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr(
			"valor inválido para o id da partida",
			400,
			"HANDLER:FOOTBALL:GET_MATCH:INVALID_MATCH_ID",
		))

		return
	}

	output, err := h.usecase.GetMatchSummary(int64(matchID))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		errors.HandleHttpError(w, err)

		return
	}
}

func (h *Handler) GetMatch(w http.ResponseWriter, r *http.Request) {
	matchID, err := strconv.ParseUint(r.PathValue("matchID"), 10, 64)
	if err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr(
			"valor inválido para o id da partida",
			400,
			"HANDLER:FOOTBALL:GET_MATCH:INVALID_MATCH_ID",
		))

		return
	}

	match, err := h.usecase.FindOneByID(int64(matchID))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func (h *Handler) FindMatchesBatch(w http.ResponseWriter, r *http.Request) {
	var IDs []int64

	if err := json.NewDecoder(r.Body).Decode(&IDs); err != nil {
		err := errors.NewHTTPErr(
			"corpo de requisição inválido",
			400,
			"HANDLER:FOOTBALL:FIND_MATCHES_BATCH:INVALID_BODY",
		)

		errors.HandleHttpError(w, err)

		return
	}

	matches, err := h.usecase.FindBatchByID(IDs)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func (h *Handler) ListMatches(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.NewPageWithRequest(r)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	team, err := strconv.ParseInt(r.URL.Query().Get("team"), 10, 64)
	if err != nil {
		team = 0
	}

	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		year = 0
	}

	league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		league = 0
	}

	status := r.URL.Query().Get("status")

	matches, err := h.usecase.List(ListMatchOptions{Team: team, Year: year, Status: status, League: league}, page.Size, page.Index)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(matches); err != nil {
		errors.HandleHttpError(w, err)
		return
	}
}
