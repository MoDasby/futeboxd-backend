package handler

import (
	"encoding/json"
	"net/http"

	"github.com/modasby/futeboxd-api/services/football/internal/errors"
	"github.com/modasby/futeboxd-api/services/football/internal/usecase"
)

type FootballHandler struct {
	listTeamsUseCase       *usecase.ListTeamsUseCase
	getTeamScheduleUseCase *usecase.GetTeamScheduleUseCase
	getMatchUseCase        *usecase.GetMatchUseCase
	getTeamByIdUseCase     *usecase.GetTeamByIdUseCase
}

func NewFootballHandler(
	listTeamsUseCase *usecase.ListTeamsUseCase,
	getTeamUseCase *usecase.GetTeamScheduleUseCase,
	getMatchUseCase *usecase.GetMatchUseCase,
	getTeamByIdUseCase *usecase.GetTeamByIdUseCase,
) *FootballHandler {
	return &FootballHandler{
		listTeamsUseCase:       listTeamsUseCase,
		getTeamScheduleUseCase: getTeamUseCase,
		getMatchUseCase:        getMatchUseCase,
		getTeamByIdUseCase:     getTeamByIdUseCase,
	}
}

func (h *FootballHandler) FindTeamByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	teams, err := h.listTeamsUseCase.Execute(name)
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

	schedule, err := h.getTeamScheduleUseCase.Execute("all", team, season)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

func (h *FootballHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamID")

	team, err := h.getTeamByIdUseCase.Execute(teamID)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *FootballHandler) GetMatchSummary(w http.ResponseWriter, r *http.Request) {
	matchID := r.PathValue("matchID")

	match, err := h.getMatchUseCase.Execute(matchID)
	if err != nil {
		errors.HandleHttpError(w, err)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}
