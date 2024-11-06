package repository

import (
	"database/sql"
	"strings"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/errors"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) domain.TeamRepository {
	return &teamRepository{db: db}
}

// AddTeam insere um novo time na tabela de teams.
func (repo *teamRepository) Create(team *domain.Team) error {
	query := `
		INSERT INTO teams (id, name, abbreviation, color, logo)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := repo.db.Exec(query, team.ID, team.Name, team.Abbreviation, team.Color, team.Logo)
	return err
}

// FindTeamById busca um time pelo ID.
func (repo *teamRepository) FindOneById(ID string) (*domain.Team, error) {
	query := `
		SELECT t.id, t.name, t.abbreviation, t.color, t.logo
		FROM teams t
		WHERE t.id = $1
	`

	row := repo.db.QueryRow(query, strings.TrimSpace(ID))

	var teamName, teamAbbreviation, teamColor, teamLogo sql.NullString
	var teamID sql.NullInt64

	if err := row.Scan(&teamID, &teamName, &teamAbbreviation, &teamColor, &teamLogo); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr("time não encontrado", 404, "REPOSITORY:TEAM:FIND_BY_ID:NOT_FOUND")
		}
		return nil, err
	}

	team := &domain.Team{
		ID:           teamID.Int64,
		Name:         teamName.String,
		Abbreviation: teamAbbreviation.String,
		Color:        teamColor.String,
		Logo:         teamLogo.String,
	}

	return team, nil
}

func (repo *teamRepository) FindAll(name string) ([]*domain.Team, error) {
	query := `
		SELECT t.id, t.name, t.abbreviation, t.color, t.logo
		FROM teams t 
		WHERE LOWER(t.name) LIKE $1
	`

	rows, err := repo.db.Query(query, strings.ToLower("%"+name+"%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := make([]*domain.Team, 0)

	for rows.Next() {
		var name, abbreviation, color, logo sql.NullString
		var id sql.NullInt64

		if err := rows.Scan(&id, &name, &abbreviation, &color, &logo); err != nil {
			return nil, err
		}

		team := domain.NewTeam(id.Int64, name.String, abbreviation.String, color.String, logo.String)

		teams = append(teams, team)
	}

	return teams, nil
}
