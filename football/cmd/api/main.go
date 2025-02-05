package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/football/database"
	"github.com/modasby/futeboxd-api/services/football/internal/match"
	"github.com/modasby/futeboxd-api/services/football/internal/news"
	"github.com/modasby/futeboxd-api/services/football/internal/team"
)

func main() {
	db := database.InitDatabase()

	teamRepository := team.NewTeamRepository(db)
	matchRepository := match.NewMatchRepository(db)
	newsRepository := news.NewNewsRepo(db)

	matchUC := match.NewMatchUsecases(matchRepository)

	matchHandler := match.NewHandler(matchUC)
	teamHandler := team.NewHandler(teamRepository)
	newsHandler := news.NewNewsHandler(newsRepository)

	r := http.NewServeMux()

	matchHandler.RegisterRoutes(r)
	teamHandler.RegisterRoutes(r)
	newsHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	log.Printf("Iniciando servidor na porta: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
