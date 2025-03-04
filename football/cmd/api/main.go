package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/football/database"
	"github.com/modasby/futeboxd-api/services/football/internal/match"
	"github.com/modasby/futeboxd-api/services/football/internal/news"
	"github.com/modasby/futeboxd-api/services/football/internal/team"
	"github.com/modasby/futeboxd-api/services/football/pkg/log"
	"github.com/modasby/futeboxd-api/services/football/pkg/middleware"
)

func main() {
	ctx := context.Background()

	if err := log.InitLogger(); err != nil {
		panic(err)
	}

	db, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}

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

	slog.Debug(fmt.Sprintf("Iniciando servidor na porta: %s", port))

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), middleware.LoggerMiddleware(ctx, r.ServeHTTP)); err != nil {
		slog.Error("Erro no listen and serve", "originalError", err)
	}
}
