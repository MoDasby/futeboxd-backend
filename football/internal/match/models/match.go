package models

import (
	"time"

	"github.com/modasby/futeboxd-api/services/football/internal/team"
)

type Competitor struct {
	HomeAway string    `json:"homeAway"`
	Winner   bool      `json:"winner"`
	Score    int32     `json:"score"`
	Team     team.Team `json:"team"`
}

type Match struct {
	ID             int64      `json:"id"`
	Venue          string     `json:"venue"`
	Datetime       time.Time  `json:"datetime"`
	Note           string     `json:"note"`
	Completed      bool       `json:"completed"`
	StatusName     string     `json:"status_name"`
	League         League     `json:"league"`
	HomeCompetitor Competitor `json:"home_competitor"`
	AwayCompetitor Competitor `json:"away_competitor"`
}

func NewMatch(
	ID int64,
	venue string,
	datetime time.Time,
	note string,
	homeCompetitor Competitor,
	awayCompetitor Competitor,
	completed bool,
	statusName string,
	leagueID int,
	leagueName string,
	leagueLogo string,
) *Match {
	return &Match{
		ID:             ID,
		Venue:          venue,
		Datetime:       datetime,
		Note:           note,
		HomeCompetitor: homeCompetitor,
		AwayCompetitor: awayCompetitor,
		Completed:      completed,
		StatusName:     statusName,
		League: League{
			ID:   leagueID,
			Name: leagueName,
			Logo: leagueLogo,
		},
	}
}

func NewCompetitor(homeAway string, winner bool, score int32, team team.Team) *Competitor {
	return &Competitor{
		HomeAway: homeAway, Winner: winner, Score: score, Team: team,
	}
}
