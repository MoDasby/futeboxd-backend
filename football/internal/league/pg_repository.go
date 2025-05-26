package league

import (
	"context"
	"database/sql"
)

type LeagueRepository struct {
	db *sql.DB
}

func NewLeagueRepository(db *sql.DB) Repository {
	return &LeagueRepository{
		db: db,
	}
}

func (repo *LeagueRepository) ListLeagues(ctx context.Context) ([]League, error) {
	query := `
		SELECT id, name, logo FROM leagues
	`

	result := make([]League, 0)

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var league League
		if err := rows.Scan(&league.ID, &league.Name, &league.Logo); err != nil {
			return nil, err
		}

		result = append(result, league)
	}

	return result, nil
}
