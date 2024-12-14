package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/errors"
)

type matchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) domain.MatchRepository {
	return &matchRepository{db: db}
}

func (repo *matchRepository) Create(match *domain.Match, summary *domain.MatchSummary) error {
	query := `
		INSERT INTO matches (
			id, match_date, venue, home_team_score, home_team_id,
			away_team_score, away_team_id, 
			note, completed, status_name, competition_name, events
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	events, err := json.Marshal(summary.Events)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(
		query,
		match.ID, match.Date, match.Venue, match.HomeCompetitor.Score,
		match.HomeCompetitor.Team.ID,
		match.AwayCompetitor.Score, match.AwayCompetitor.Team.ID,
		match.Note, match.Completed, match.StatusName, match.CompetitionName, events,
	)

	return err
}

func (repo *matchRepository) FindOneByID(matchID int64) (*domain.Match, error) {

	query := `
		SELECT m.id, match_date, venue, 
		home_team_score, home_team_id, ht.name AS home_team_name, 
		ht.abbreviation AS home_team_abbreviation, ht.color AS home_team_color, ht.logo AS home_team_logo,
		away_team_score, away_team_id, at.name AS away_team_name, 
		at.abbreviation AS away_team_abbreviation, at.color AS away_team_color, at.logo AS away_team_logo,
		note, completed, status_name, competition_name, 
		(home_team_score > away_team_score) AS home_team_winner,
		(away_team_score > home_team_score) AS away_team_winner
		FROM matches m
		LEFT JOIN teams ht ON ht.id = m.home_team_id
		LEFT JOIN teams at ON at.id = m.away_team_id
		WHERE m.id = $1
	`

	rows := repo.db.QueryRow(query, matchID)

	var match domain.Match

	if err := rows.Scan(
		&match.ID, &match.Date, &match.Venue, &match.HomeCompetitor.Score,
		&match.HomeCompetitor.Team.ID, &match.HomeCompetitor.Team.Name, &match.HomeCompetitor.Team.Abbreviation,
		&match.HomeCompetitor.Team.Color, &match.HomeCompetitor.Team.Logo,
		&match.AwayCompetitor.Score, &match.AwayCompetitor.Team.ID, &match.AwayCompetitor.Team.Name,
		&match.AwayCompetitor.Team.Abbreviation, &match.AwayCompetitor.Team.Color, &match.AwayCompetitor.Team.Logo,
		&match.Note, &match.Completed, &match.StatusName, &match.CompetitionName,
		&match.HomeCompetitor.Winner, &match.AwayCompetitor.Winner,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:FIND_BY_ID:NOT_FOUND")
		}

		return nil, err
	}

	return &match, nil
}

func (repo *matchRepository) FindBatchByID(ids []int64) ([]domain.Match, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)) // slice de argumentos

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	placeholderStr := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(`
		SELECT 
			m.id, match_date, venue, 
			home_team_score, home_team_id, ht.name AS home_team_name, 
			ht.abbreviation AS home_team_abbreviation, ht.color AS home_team_color, ht.logo AS home_team_logo,
			away_team_score, away_team_id, at.name AS away_team_name, 
			at.abbreviation AS away_team_abbreviation, at.color AS away_team_color, at.logo AS away_team_logo,
			note, completed, status_name, competition_name, 
			(home_team_score > away_team_score) AS home_team_winner,
			(away_team_score > home_team_score) AS away_team_winner
		FROM matches m
		LEFT JOIN teams ht ON ht.id = m.home_team_id
		LEFT JOIN teams at ON at.id = m.away_team_id
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

		if err := rows.Scan(
			&match.ID, &match.Date, &match.Venue, &match.HomeCompetitor.Score,
			&match.HomeCompetitor.Team.ID, &match.HomeCompetitor.Team.Name, &match.HomeCompetitor.Team.Abbreviation,
			&match.HomeCompetitor.Team.Color, &match.HomeCompetitor.Team.Logo,
			&match.AwayCompetitor.Score, &match.AwayCompetitor.Team.ID, &match.AwayCompetitor.Team.Name,
			&match.AwayCompetitor.Team.Abbreviation, &match.AwayCompetitor.Team.Color, &match.AwayCompetitor.Team.Logo,
			&match.Note, &match.Completed, &match.StatusName, &match.CompetitionName,
			&match.HomeCompetitor.Winner, &match.AwayCompetitor.Winner,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:FIND_BY_ID:NOT_FOUND")
			}

			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (repo *matchRepository) GetMatchSummary(matchID int64) (*domain.MatchSummary, error) {
	query := `
		SELECT events 
		FROM matches
		WHERE id = $1
	`

	row := repo.db.QueryRow(query, matchID)

	var output domain.MatchSummary

	var events sql.NullString

	if err := row.Scan(&events); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:GET_MATCH_SUMMARY:NOT_FOUND")
		}

		return nil, err
	}

	if events.Valid {
		if err := json.Unmarshal([]byte(events.String), &output.Events); err != nil {
			return nil, err
		}
	}

	return &output, nil
}

func (repo *matchRepository) ListByYear(teamID int64, year, pageSize, pageIndex int) ([]domain.Match, error) {
	query := `
		SELECT 
			m.id, match_date, venue, 
			home_team_score, home_team_id, ht.name AS home_team_name, 
			ht.abbreviation AS home_team_abbreviation, ht.color AS home_team_color, ht.logo AS home_team_logo,
			away_team_score, away_team_id, at.name AS away_team_name, 
			at.abbreviation AS away_team_abbreviation, at.color AS away_team_color, at.logo AS away_team_logo,
			note, completed, status_name, competition_name, 
			(home_team_score > away_team_score) AS home_team_winner,
			(away_team_score > home_team_score) AS away_team_winner
		FROM matches m
		LEFT JOIN teams ht ON ht.id = m.home_team_id
		LEFT JOIN teams at ON at.id = m.away_team_id
		WHERE (m.home_team_id = $1 or m.away_team_id = $1) and EXTRACT(YEAR FROM m.match_date) = $2
		ORDER BY m.match_date DESC
		LIMIT $3
		OFFSET ($4 - 1) * $3
	`

	rows, err := repo.db.Query(query, teamID, year, pageSize, pageIndex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]domain.Match, 0)

	for rows.Next() {
		var match domain.Match

		if err := rows.Scan(
			&match.ID, &match.Date, &match.Venue, &match.HomeCompetitor.Score,
			&match.HomeCompetitor.Team.ID, &match.HomeCompetitor.Team.Name, &match.HomeCompetitor.Team.Abbreviation,
			&match.HomeCompetitor.Team.Color, &match.HomeCompetitor.Team.Logo,
			&match.AwayCompetitor.Score, &match.AwayCompetitor.Team.ID, &match.AwayCompetitor.Team.Name,
			&match.AwayCompetitor.Team.Abbreviation, &match.AwayCompetitor.Team.Color, &match.AwayCompetitor.Team.Logo,
			&match.Note, &match.Completed, &match.StatusName, &match.CompetitionName,
			&match.HomeCompetitor.Winner, &match.AwayCompetitor.Winner,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:FIND_BY_ID:NOT_FOUND")
			}

			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (repo *matchRepository) FindLiveMatches(pageSize, pageIndex int) ([]domain.Match, error) {
	query := `
		SELECT 
			m.id, match_date, venue, 
			home_team_score, home_team_id, ht.name AS home_team_name, 
			ht.abbreviation AS home_team_abbreviation, ht.color AS home_team_color, ht.logo AS home_team_logo,
			away_team_score, away_team_id, at.name AS away_team_name, 
			at.abbreviation AS away_team_abbreviation, at.color AS away_team_color, at.logo AS away_team_logo,
			note, completed, status_name, competition_name, 
			(home_team_score > away_team_score) AS home_team_winner,
			(away_team_score > home_team_score) AS away_team_winner
		FROM matches m
		LEFT JOIN teams ht ON ht.id = m.home_team_id
		LEFT JOIN teams at ON at.id = m.away_team_id
		WHERE NOT m.completed AND (m.status_name = 'STATUS_FIRST_HALF' OR m.status_name = 'STATUS_SECOND_HALF' OR m.status_name = 'STATUS_HALFTIME')
		ORDER BY m.match_date
		LIMIT $1
		OFFSET ($2 - 1) * $1
	`

	rows, err := repo.db.Query(query, pageSize, pageIndex)
	if err != nil {
		return nil, err
	}

	matches := make([]domain.Match, 0)

	for rows.Next() {
		var match domain.Match

		if err := rows.Scan(
			&match.ID, &match.Date, &match.Venue, &match.HomeCompetitor.Score,
			&match.HomeCompetitor.Team.ID, &match.HomeCompetitor.Team.Name, &match.HomeCompetitor.Team.Abbreviation,
			&match.HomeCompetitor.Team.Color, &match.HomeCompetitor.Team.Logo,
			&match.AwayCompetitor.Score, &match.AwayCompetitor.Team.ID, &match.AwayCompetitor.Team.Name,
			&match.AwayCompetitor.Team.Abbreviation, &match.AwayCompetitor.Team.Color, &match.AwayCompetitor.Team.Logo,
			&match.Note, &match.Completed, &match.StatusName, &match.CompetitionName,
			&match.HomeCompetitor.Winner, &match.AwayCompetitor.Winner,
		); err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}
