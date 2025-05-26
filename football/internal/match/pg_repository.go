package match

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modasby/futeboxd-api/services/football/internal/match/models"
	"github.com/modasby/futeboxd-api/services/football/pkg/errors"
	"github.com/modasby/futeboxd-api/services/football/pkg/sq"
)

type matchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) Repository {
	return &matchRepository{db: db}
}

func (repo *matchRepository) FindOneByID(matchID int64) (*models.Match, error) {

	query := `
		SELECT m.id, match_date, venue, 
		home_team_score, home_team_id, ht.name AS home_team_name, 
		ht.abbreviation AS home_team_abbreviation, ht.color AS home_team_color, ht.logo AS home_team_logo,
		away_team_score, away_team_id, at.name AS away_team_name, 
		at.abbreviation AS away_team_abbreviation, at.color AS away_team_color, at.logo AS away_team_logo,
		note, completed, status_name, l.id, l.name, l.logo, 
		(home_team_score > away_team_score) AS home_team_winner,
		(away_team_score > home_team_score) AS away_team_winner
		FROM matches m
		LEFT JOIN teams ht ON ht.id = m.home_team_id
		LEFT JOIN teams at ON at.id = m.away_team_id
		LEFT JOIN leagues l ON l.id = m.league_id
		WHERE m.id = $1
	`

	rows := repo.db.QueryRow(query, matchID)

	var match models.Match
	var note sql.NullString

	if err := rows.Scan(
		&match.ID, &match.Datetime, &match.Venue, &match.HomeCompetitor.Score,
		&match.HomeCompetitor.Team.ID, &match.HomeCompetitor.Team.Name, &match.HomeCompetitor.Team.Abbreviation,
		&match.HomeCompetitor.Team.Color, &match.HomeCompetitor.Team.Logo,
		&match.AwayCompetitor.Score, &match.AwayCompetitor.Team.ID, &match.AwayCompetitor.Team.Name,
		&match.AwayCompetitor.Team.Abbreviation, &match.AwayCompetitor.Team.Color, &match.AwayCompetitor.Team.Logo,
		&note, &match.Completed, &match.StatusName, &match.League.ID, &match.League.Name,
		&match.League.Logo, &match.HomeCompetitor.Winner, &match.AwayCompetitor.Winner,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:FIND_BY_ID:NOT_FOUND")
		}

		return nil, err
	}

	if note.Valid {
		match.Note = note.String
	}

	return &match, nil
}

func (repo *matchRepository) FindBatchByID(ids []int64) ([]models.Match, error) {
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
			note, completed, status_name, l.id, l.name, l.logo,
			(home_team_score > away_team_score) AS home_team_winner,
			(away_team_score > home_team_score) AS away_team_winner
		FROM matches m
		LEFT JOIN teams ht ON ht.id = m.home_team_id
		LEFT JOIN teams at ON at.id = m.away_team_id
		LEFT JOIN leagues l ON l.id = m.league_id
		WHERE m.id IN (%s)
	`, placeholderStr)

	rows, err := repo.db.Query(query, args...) // Use args com o operador variadic
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]models.Match, 0)

	for rows.Next() {
		var match models.Match
		var note sql.NullString

		if err := rows.Scan(
			&match.ID, &match.Datetime, &match.Venue, &match.HomeCompetitor.Score,
			&match.HomeCompetitor.Team.ID, &match.HomeCompetitor.Team.Name, &match.HomeCompetitor.Team.Abbreviation,
			&match.HomeCompetitor.Team.Color, &match.HomeCompetitor.Team.Logo,
			&match.AwayCompetitor.Score, &match.AwayCompetitor.Team.ID, &match.AwayCompetitor.Team.Name,
			&match.AwayCompetitor.Team.Abbreviation, &match.AwayCompetitor.Team.Color, &match.AwayCompetitor.Team.Logo,
			&note, &match.Completed, &match.StatusName, &match.League.ID, &match.League.Name,
			&match.League.Logo, &match.HomeCompetitor.Winner, &match.AwayCompetitor.Winner,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:FIND_BY_ID:NOT_FOUND")
			}

			return nil, err
		}

		if note.Valid {
			match.Note = note.String
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (repo *matchRepository) GetMatchSummary(matchID int64) (*models.MatchSummary, error) {
	query := `
		SELECT events 
		FROM matches
		WHERE id = $1
	`

	row := repo.db.QueryRow(query, matchID)

	var summary models.MatchSummary
	summary.Events = make([]models.Event, 0)

	var eventsStr sql.NullString

	if err := row.Scan(&eventsStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr("partida não encontrada", 404, "REPOSITORY:MATCH:GET_MATCH_SUMMARY:NOT_FOUND")
		}

		return nil, err
	}

	if eventsStr.Valid {
		if err := json.Unmarshal([]byte(eventsStr.String), &summary.Events); err != nil {
			return nil, err
		}
	}

	return &summary, nil
}

func (repo *matchRepository) List(filters ListMatchOptions, pageSize, pageIndex int) ([]models.Match, error) {
	query := sq.Query.Select(`
			m.id, m.match_date AS match_date, venue, 
			home_team_score, home_team_id, ht.name AS home_team_name, 
			ht.abbreviation AS home_team_abbreviation, ht.color AS home_team_color, ht.logo AS home_team_logo,
			away_team_score, away_team_id, at.name AS away_team_name, 
			at.abbreviation AS away_team_abbreviation, at.color AS away_team_color, at.logo AS away_team_logo,
			note, completed, status_name, l.id, l.name, l.logo,
			(home_team_score > away_team_score) AS home_team_winner,
			(away_team_score > home_team_score) AS away_team_winner
		`).
		From(`matches m`).
		LeftJoin(`teams ht ON ht.id = m.home_team_id`).
		LeftJoin(`teams at ON at.id = m.away_team_id`).
		LeftJoin(`leagues l ON l.id = m.league_id`)

	if filters.League > 0 {
		query = query.Where("league_id = ?", filters.League)
	}

	if filters.Team > 0 {
		query = query.Where("(home_team_id = ? OR away_team_id = ?)", filters.Team, filters.Team)
	}

	if filters.Year > 0 {
		query = query.Where("EXTRACT(YEAR FROM m.match_date) = ?", filters.Year)
	}

	if filters.Status == "completed" {
		query = query.
			Where(sq.NotEq{"m.status_name": "STATUS_FIRST_HALF"}).
			Where(sq.NotEq{"m.status_name": "STATUS_SCHEDULED"}).
			Where(sq.NotEq{"m.status_name": "STATUS_SECOND_HALF"}).
			Where(sq.NotEq{"m.status_name": "STATUS_HALFTIME"}).
			Where(sq.NotEq{"m.status_name": "STATUS_OVERTIME"}).
			Where(sq.NotEq{"m.status_name": "STATUS_SHOOTOUT"}).
			Where(sq.NotEq{"m.status_name": "STATUS_POSTPONED"})
	}

	if filters.Status == "live" {
		query = query.
			Where(sq.Eq{"completed": false}).
			Where(sq.NotEq{"m.status_name": "STATUS_SCHEDULED"}).
			Where(sq.NotEq{"m.status_name": "STATUS_POSTPONED"}).
			Where(sq.NotEq{"m.status_name": "STATUS_ABANDONED"})
	}

	query = query.
		Limit(uint64(pageSize)).
		Offset(uint64(pageIndex - 1)).
		OrderBy("m.match_date DESC")

	rows, err := query.RunWith(repo.db).Query()
	if err != nil {
		return nil, err
	}

	matches := make([]models.Match, 0)

	for rows.Next() {
		var match models.Match
		var note sql.NullString

		if err := rows.Scan(
			&match.ID, &match.Datetime, &match.Venue, &match.HomeCompetitor.Score,
			&match.HomeCompetitor.Team.ID, &match.HomeCompetitor.Team.Name, &match.HomeCompetitor.Team.Abbreviation,
			&match.HomeCompetitor.Team.Color, &match.HomeCompetitor.Team.Logo,
			&match.AwayCompetitor.Score, &match.AwayCompetitor.Team.ID, &match.AwayCompetitor.Team.Name,
			&match.AwayCompetitor.Team.Abbreviation, &match.AwayCompetitor.Team.Color, &match.AwayCompetitor.Team.Logo,
			&note, &match.Completed, &match.StatusName, &match.League.ID, &match.League.Name,
			&match.League.Logo, &match.HomeCompetitor.Winner, &match.AwayCompetitor.Winner,
		); err != nil {
			return nil, err
		}

		if note.Valid {
			match.Note = note.String
		}

		matches = append(matches, match)
	}

	return matches, nil
}
