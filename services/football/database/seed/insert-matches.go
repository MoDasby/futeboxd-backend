package main

import (
	"log"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/database"
	"github.com/modasby/futeboxd-api/services/football/internal/repository"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
	"github.com/modasby/futeboxd-api/services/football/internal/usecase"
)

func InsertMatches() {
	db := database.InitDatabase()

	teamRepository := repository.NewTeamRepository(db)
	matchRepository := repository.NewMatchRepository(db)

	espnService := service.NewEspnService()

	getSchedule := usecase.NewGetTeamScheduleUseCase(espnService)
	getMatch := usecase.NewGetMatchUseCase(espnService, matchRepository)

	teams, _ := teamRepository.FindAll("")

	for _, team := range teams {
		var temp int64
		for temp = 2006; temp < 2025; temp++ {
			schedule, err := getSchedule.Execute("all", strconv.FormatInt(team.ID, 10), strconv.FormatInt(temp, 10))
			if err != nil {
				panic(err)
			}

			for _, match := range schedule.Matches {
				match, err := getMatch.Execute(strconv.FormatInt(match.ID, 10))
				if err != nil {
					panic(err)
				}

				log.Printf("iserida partida %d", match.ID)
			}

			log.Printf("inserido agenda do time %s da temporada %d\n", team.Name, temp)
		}
	}
}
