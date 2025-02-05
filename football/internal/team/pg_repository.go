package team

import (
	"database/sql"
	"strings"

	"github.com/modasby/futeboxd-api/services/football/pkg/errors"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) Repository {
	return &teamRepository{db: db}
}

func (repo *teamRepository) FindOneById(ID int64) (*Team, error) {
	query := `
		SELECT t.id, t.name, t.abbreviation, t.color, t.logo
		FROM teams t
		WHERE t.id = $1
	`

	row := repo.db.QueryRow(query, ID)

	var teamName, teamAbbreviation, teamColor, teamLogo sql.NullString
	var teamID sql.NullInt64

	if err := row.Scan(&teamID, &teamName, &teamAbbreviation, &teamColor, &teamLogo); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr("time não encontrado", 404, "REPOSITORY:TEAM:FIND_BY_ID:NOT_FOUND")
		}
		return nil, err
	}

	team := &Team{
		ID:           teamID.Int64,
		Name:         teamName.String,
		Abbreviation: teamAbbreviation.String,
		Color:        teamColor.String,
		Logo:         teamLogo.String,
	}

	return team, nil
}

func (repo *teamRepository) FindByName(name string, pageSize, pageIndex int) ([]Team, error) {
	query := `
	SELECT t.id, t.name, t.abbreviation, t.color, t.logo, t.mandatory
	FROM teams t
	WHERE unaccent(t.name) ILIKE unaccent($1) AND t.mandatory
	ORDER BY t.name
	LIMIT $2
	OFFSET (($3 - 1) * $2)
	`

	rows, err := repo.db.Query(query, strings.ToLower("%"+name+"%"), pageSize, pageIndex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := make([]Team, 0)

	for rows.Next() {
		var name, abbreviation, color, logo sql.NullString
		var id sql.NullInt64
		var mandatory sql.NullBool

		if err := rows.Scan(&id, &name, &abbreviation, &color, &logo, &mandatory); err != nil {
			return nil, err
		}

		team := NewTeam(id.Int64, name.String, abbreviation.String, color.String, logo.String)

		team.Mandatory = mandatory.Bool

		teams = append(teams, *team)
	}

	return teams, nil
}
