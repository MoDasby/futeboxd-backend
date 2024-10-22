package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/football/database"
	"github.com/modasby/futeboxd-api/services/football/internal/handler"
	"github.com/modasby/futeboxd-api/services/football/internal/repository"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
	"github.com/modasby/futeboxd-api/services/football/internal/usecase"
)

func main() {
	db := database.InitDatabase()

	espnService := service.NewEspnService()

	teamRepository := repository.NewTeamRepository(db)
	matchRepository := repository.NewMatchRepository(db)

	listTeamsUseCase := usecase.NewListTeamsUseCase(teamRepository)
	getTeamScheduleUseCase := usecase.NewGetTeamScheduleUseCase(espnService)
	getMatchUseCase := usecase.NewGetMatchUseCase(espnService, matchRepository)
	getTeamByIdUseCase := usecase.NewGetTeamByIdUseCase(teamRepository)

	footballHandler := handler.NewFootballHandler(
		listTeamsUseCase,
		getTeamScheduleUseCase,
		getMatchUseCase,
		getTeamByIdUseCase,
	)

	r := http.NewServeMux()

	r.HandleFunc("/football/team", footballHandler.FindTeamByName)
	r.HandleFunc("/football/team/{teamID}", footballHandler.GetTeamByID)
	r.HandleFunc("/football/team/{team}/schedule", footballHandler.ListTeamSchedule)
	r.HandleFunc("/football/match/{matchID}", footballHandler.GetMatchSummary)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
