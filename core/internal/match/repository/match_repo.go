package repository

import (
	"context"
	"database/sql"

	"github.com/modasby/futeboxd-backend/core/internal/domain"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type matchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) domain.MatchRepository {
	return &matchRepository{db: db}
}

func (repo *matchRepository) ListPopularMatches(ctx context.Context, page *pagination.Page) ([]int64, error) {
	query := `
		SELECT 
		    r.match_id
		FROM reviews r
		WHERE r.created_at >= NOW() - INTERVAL '7 days'
		GROUP BY r.match_id
		ORDER BY (
    		LEAST(COUNT(r.match_id), 100)
		) + (
			ROUND(AVG(r.rate), 1)
		) DESC
		LIMIT $1
		OFFSET ($2 - 1) * $1
	`

	var matches []int64

	rows, err := repo.db.QueryContext(ctx, query, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		matches = append(matches, id)
	}

	return matches, nil
}
