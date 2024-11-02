package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	errorsTypes "github.com/modasby/futeboxd-api/services/football/internal/errors"
)

type matchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) domain.MatchRepository {
	return &matchRepository{db: db}
}

func (repo *matchRepository) AddMatch(match *domain.Match) error {
	query := `
		INSERT INTO matches (
			id, match_date, venue_name, venue_city, home_team_score, home_team_id, home_team_rosters,
			away_team_score, away_team_id, away_team_rosters, 
			note, completed, status_name, competition_name, events
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	events, err := json.Marshal(match.Events)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(
		query,
		match.ID, match.Date, match.Venue.Name, match.Venue.City, match.HomeCompetitor.Score,
		match.HomeCompetitor.Team.ID, match.HomeCompetitor.Roster,
		match.AwayCompetitor.Score, match.AwayCompetitor.Team.ID, match.AwayCompetitor.Roster,
		match.Note, match.Completed, match.StatusName, match.CompetitionName, events,
	)

	return err
}

func (repo *matchRepository) FindMatchByID(matchID string) (*domain.Match, error) {

	query := `
		SELECT id, match_date, venue_name, venue_city, home_team_score, home_team_id, home_team_rosters,
			away_team_score, away_team_id, away_team_rosters, 
			note, completed, status_name, competition_name, events
		FROM matches m
		WHERE m.id = $1
	`

	rows := repo.db.QueryRow(query, matchID)

	var match domain.Match

	var homeTeamRoster, awayTeamRoster, events []byte

	if err := rows.Scan(
		&match.ID, &match.Date, &match.Venue.Name, &match.Venue.City, &match.HomeCompetitor.Score,
		&match.HomeCompetitor.Team.ID, &homeTeamRoster, &match.AwayCompetitor.Score,
		&match.AwayCompetitor.Team.ID, &awayTeamRoster, &match.Note, &match.Completed, &match.StatusName,
		&match.CompetitionName, &events,
	); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:FIND_BY_ID:NOT_FOUND")
		}

		return nil, err
	}

	if err := json.Unmarshal(homeTeamRoster, &match.HomeCompetitor.Roster); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(awayTeamRoster, &match.AwayCompetitor.Roster); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(events, &match.Events); err != nil {
		return nil, err
	}

	return &match, nil
}

func (repo *matchRepository) FindMatchBatch(ids []int64) ([]domain.Match, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)) // slice de argumentos

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	placeholderStr := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(`
		SELECT m.id, m.match_date, m.venue, m.competitors, m.note, m.completed, m.status_name,
		m.competition_name, m.events
		FROM matches m
		WHERE m.id IN (%s)
	`, placeholderStr)

	rows, err := repo.db.Query(query, args...) // Use args com o operador variadic
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]domain.Match, 0)

	for rows.Next() {
		var match domain.Match
		var venue, competitors, events []byte

		if err := rows.Scan(
			&match.ID,
			&match.Date,
			&venue,
			&competitors,
			&match.Note,
			&match.Completed,
			&match.StatusName,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(venue, &match.Venue); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(events, &match.Events); err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}
