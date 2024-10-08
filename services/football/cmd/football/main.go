package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/modasby/futeboxd-api/services/football/internal/handler"
	"github.com/modasby/futeboxd-api/services/football/internal/repository"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
	"github.com/modasby/futeboxd-api/services/football/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable port=5433",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

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
