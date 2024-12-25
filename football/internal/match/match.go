package match

import "github.com/modasby/futeboxd-api/services/football/internal/team"

type Competitor struct {
	HomeAway string    `json:"homeAway"`
	Winner   bool      `json:"winner"`
	Score    int32     `json:"score"`
	Team     team.Team `json:"team"`
}

type Match struct {
	ID              int64      `json:"id"`
	Venue           string     `json:"venue"`
	Date            string     `json:"date"`
	Note            string     `json:"note"`
	Completed       bool       `json:"completed"`
	StatusName      string     `json:"status_name"`
	CompetitionName string     `json:"competition_name"`
	HomeCompetitor  Competitor `json:"home_competitor"`
	AwayCompetitor  Competitor `json:"away_competitor"`
}

func NewMatch(
	ID int64,
	venue string,
	date string,
	note string,
	homeCompetitor Competitor,
	awayCompetitor Competitor,
	completed bool,
	statusName string,
	competitionName string,
) *Match {
	return &Match{
		ID:              ID,
		Venue:           venue,
		Date:            date,
		Note:            note,
		HomeCompetitor:  homeCompetitor,
		AwayCompetitor:  awayCompetitor,
		Completed:       completed,
		StatusName:      statusName,
		CompetitionName: competitionName,
	}
}

func NewCompetitor(homeAway string, winner bool, score int32, team team.Team) *Competitor {
	return &Competitor{
		HomeAway: homeAway, Winner: winner, Score: score, Team: team,
	}
}
