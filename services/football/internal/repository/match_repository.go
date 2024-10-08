package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

	errorsTypes "github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/football/internal/domain"
)

type matchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) domain.MatchRepository {
	return &matchRepository{db: db}
}

func (repo *matchRepository) AddMatch(match *domain.Match) error {
	query := `
		INSERT INTO matches (id, match_date, venue, competitors, note, completed, status_name, competition_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	venue, err := json.Marshal(match.Venue)
	if err != nil {
		return err
	}

	competitors, err := json.Marshal(match.Competitors)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(query, match.ID, match.Date, venue, competitors, match.Note, match.Completed, match.StatusName, match.CompetitionName)

	return err
}

func (repo *matchRepository) FindMatchByID(matchID string) (*domain.Match, error) {

	query := `
		SELECT m.id, m.match_date, m.venue, m.competitors, m.note, m.completed, m.status_name,
		m.competition_name
		FROM matches m
		WHERE m.id = $1
	`

	rows := repo.db.QueryRow(query, matchID)

	var match *domain.Match

	var venue, competitors []byte
	var date, statusName, competitionName, note sql.NullString
	var ID sql.NullInt64
	var completed sql.NullBool

	if err := rows.Scan(&ID, &date, &venue, &competitors, &note, &completed, &statusName, &competitionName); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewErrNotFound("partida não encontrada")
		}

		return nil, err
	}

	match = &domain.Match{
		ID:              ID.Int64,
		Date:            date.String,
		Completed:       completed.Bool,
		StatusName:      statusName.String,
		CompetitionName: competitionName.String,
		Note:            note.String,
	}

	if err := json.Unmarshal(venue, &match.Venue); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(competitors, &match.Competitors); err != nil {
		return nil, err
	}

	return match, nil
}
