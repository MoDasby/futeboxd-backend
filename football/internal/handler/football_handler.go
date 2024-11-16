package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/errors"
	"github.com/modasby/futeboxd-api/services/football/internal/service/matches"
	"github.com/modasby/futeboxd-api/services/football/internal/service/teams"
)

type FootballHandler struct {
	teamService  teams.TeamService
	matchService matches.MatchService
}

func NewFootballHandler(
	teamService teams.TeamService,
	matchService matches.MatchService,
) *FootballHandler {
	return &FootballHandler{
		teamService:  teamService,
		matchService: matchService,
	}
}

func (h *FootballHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /football/team", h.FindTeamByName)
	r.HandleFunc("GET /football/team/{teamID}", h.GetTeamByID)
	r.HandleFunc("GET /football/team/{team}/schedule", h.ListTeamSchedule)
	r.HandleFunc("GET /football/match/{matchID}", h.GetMatch)
	r.HandleFunc("GET /football/match/{matchID}/summary", h.GetMatchSummary)
	r.HandleFunc("POST /football/match/batch", h.FindMatchesBatch)
}

func (h *FootballHandler) FindMatchesBatch(w http.ResponseWriter, r *http.Request) {
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

	matches, err := h.matchService.FindBatchByID(IDs)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func (h *FootballHandler) FindTeamByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	pageSize, err := strconv.ParseUint(r.URL.Query().Get("size"), 10, 64)
	if err != nil {
		pageSize = 50
	}
	pageIndex, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		pageIndex = 1
	}

	teams, err := h.teamService.FindAll(name, int(pageSize), int(pageIndex))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *FootballHandler) ListTeamSchedule(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	season := r.URL.Query().Get("season")

	schedule, err := h.teamService.ListSchedule(team, season)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

func (h *FootballHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.ParseUint(r.PathValue("teamID"), 10, 64)
	if err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr(
			"valor inválido para o id do time",
			400,
			"HANDLER:FOOTBALL:GET_TEAM:INVALID_TEAM_ID",
		))

		return
	}

	team, err := h.teamService.FindOneById(int64(teamID))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *FootballHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	matchID, err := strconv.ParseUint(r.PathValue("matchID"), 10, 64)
	if err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr(
			"valor inválido para o id da partida",
			400,
			"HANDLER:FOOTBALL:GET_MATCH:INVALID_MATCH_ID",
		))

		return
	}

	match, err := h.matchService.FindOneByID(int64(matchID))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func (h *FootballHandler) GetMatchSummary(w http.ResponseWriter, r *http.Request) {
	matchID, err := strconv.ParseUint(r.PathValue("matchID"), 10, 64)
	if err != nil {
		errors.HandleHttpError(w, errors.NewHTTPErr(
			"valor inválido para o id da partida",
			400,
			"HANDLER:FOOTBALL:GET_MATCH:INVALID_MATCH_ID",
		))

		return
	}

	output, err := h.matchService.GetMatchSummary(int64(matchID))
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
