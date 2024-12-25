package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/football/internal/team"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres dbname=futeboxd-football password=123456 host=localhost sslmode=disable port=5433"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	teamsRepo := team.NewTeamRepository(db)

	teams, err := teamsRepo.FindAll("", 2000, 1)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll("./logos", os.ModePerm); err != nil {
		panic(err)
	}

	for _, team := range teams {
		if team.Logo == "" {
			continue
		}

		file, err := os.Create(fmt.Sprintf("logos/%d.png", team.ID))
		if err != nil {
			panic(err)
		}

		resp, err := http.Get(team.Logo)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(file, resp.Body); err != nil {
			panic(err)
		}

		file.Close()
		resp.Body.Close()
	}
}
