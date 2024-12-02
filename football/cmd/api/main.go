package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/football/database"
	"github.com/modasby/futeboxd-api/services/football/internal/handler"
	"github.com/modasby/futeboxd-api/services/football/internal/repository"
	"github.com/modasby/futeboxd-api/services/football/internal/service/espn"
	"github.com/modasby/futeboxd-api/services/football/internal/service/matches"
	"github.com/modasby/futeboxd-api/services/football/internal/service/teams"
)

func main() {
	db := database.InitDatabase()

	espnService := espn.NewEspnService()

	teamRepository := repository.NewTeamRepository(db)
	matchRepository := repository.NewMatchRepository(db)

	teamService := teams.NewTeamService(teamRepository, espnService)
	matchService := matches.NewMatchService(matchRepository, espnService)

	footballHandler := handler.NewFootballHandler(
		teamService,
		matchService,
	)

	r := http.NewServeMux()

	footballHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
