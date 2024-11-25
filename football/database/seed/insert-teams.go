package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
)

var leagues = []string{
	"bra.1", "bra.2", "bra.3", "bra.copa_do_brazil", "uefa.champions", "uefa.europa",
	"fifa.world", "conmebol.libertadores", "conmebol.sudamericana", "eng.1", "eng.2",
	"esp.1", "fra.1", "ger.1", "ita.1", "por.1", "bra.camp.paulista",
	"bra.camp.carioca", "bra.camp.gaucho", "bra.camp.mineiro", "bra.copa_do_nordeste", "bra.supercopa_do_brazil",
	"uefa.euroq", "conmebol.america",
}

func InsertTeams() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
	espn_service := service.NewEspnService()

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable port=5433",
		os.Getenv("DB_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	for _, league := range leagues {
		espnTeams, err := espn_service.ListTeams(league)
		if err != nil {
			panic(err)
		}

		teams := espnTeams.Sports[0].Leagues[0].Teams

		for _, team := range teams {
			query := `
				INSERT INTO teams (id, name, abbreviation, color, logo)
				VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING
			`

			var logo = ""

			if len(team.Team.Logos) > 0 {
				logo = team.Team.Logos[0].Href
			}

			_, err := db.Exec(
				query,
				team.Team.ID,
				team.Team.Name,
				team.Team.Abbreviation,
				team.Team.Color,
				logo,
			)
			if err != nil {
				panic(err)
			}

			log.Printf("time inserido: %s", team.Team.Name)
		}
	}
}
