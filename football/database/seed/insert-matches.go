package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/repository"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
	"github.com/modasby/futeboxd-api/services/football/internal/service/mapper"
	"github.com/modasby/futeboxd-api/services/football/internal/service/matches"
	"github.com/modasby/futeboxd-api/services/football/internal/service/teams"
)

func InsertMatches() {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable port=5433",
		"postgres",
		"futeboxd-football",
		"123456",
		"localhost",
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	teamRepository := repository.NewTeamRepository(db)
	matchRepository := repository.NewMatchRepository(db)

	espnService := service.NewEspnService()

	teamService := teams.NewTeamService(teamRepository, espnService)
	matchService := matches.NewMatchService(matchRepository, espnService)

	teams, _ := teamRepository.FindAll("", 1000, 1)

	matchesFailed := make([]int64, 0)

	for _, team := range teams {
		var temp int64
		for temp = 2006; temp < 2025; temp++ {
			schedule, err := teamService.ListSchedule(strconv.FormatInt(team.ID, 10), strconv.FormatInt(temp, 10))
			if err != nil {
				panic(err)
			}

			for _, match := range schedule.Matches {
				checkTeamExists(teamRepository, espnService, match.HomeCompetitor.Team)
				checkTeamExists(teamRepository, espnService, match.AwayCompetitor.Team)

				_, err := matchService.FindOneByID(match.ID)
				if err != nil {
					matchesFailed = append(matchesFailed, match.ID)
					log.Printf("partida falhada %d", match.ID)
				} else {
					log.Printf("inserida partida %d", match.ID)
				}
			}

			log.Printf("inserido agenda do time %s da temporada %d\n", team.Name, temp)
		}
	}

	log.Println(matchesFailed)
}

func checkTeamExists(
	teamRepository domain.TeamRepository,
	espnService service.EspnService,
	team domain.Team,
) {
	_, err := teamRepository.FindOneById(team.ID)
	if err != nil {
		espnTeam, err := espnService.GetTeam(team.ID)
		if err != nil {
			panic(err)
		}

		team, err := mapper.EspnTeamToTeam(espnTeam.Team)
		if err != nil {
			log.Fatal(err)
		}

		if err := teamRepository.Create(team); err != nil {
			log.Fatal(err)
		}
	}
}
