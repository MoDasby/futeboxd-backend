package repository

import (
	"context"
	"database/sql"

	"github.com/modasby/futeboxd-backend/core/internal/match"
	"github.com/modasby/futeboxd-backend/core/internal/match/entity"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type matchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) match.Repository {
	return &matchRepository{db: db}
}

func (repo *matchRepository) ListPopularMatches(ctx context.Context, page *pagination.Page) ([]int64, error) {
	query := `
		SELECT 
		    r.match_id
		FROM reviews r
		WHERE r.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY r.match_id
		ORDER BY (
    		LEAST(COUNT(r.match_id), 100)
		) * (
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

func (repo *matchRepository) GetMatchStats(ctx context.Context, matchID int64) (*entity.MatchStats, error) {
	query := `
	SELECT 
		match_id,
		ROUND(AVG(rate), 2) as avg_rate,
		ROUND(COUNT(CASE WHEN rate = 1 THEN 1 END), 2) as count1,
		ROUND(COUNT(CASE WHEN rate = 2 THEN 1 END), 2) as count2,
		ROUND(COUNT(CASE WHEN rate = 3 THEN 1 END), 2) as count3,
		ROUND(COUNT(CASE WHEN rate = 4 THEN 1 END), 2) as count4,
		ROUND(COUNT(CASE WHEN rate = 5 THEN 1 END), 2) as count5
	FROM reviews
	WHERE match_id = $1
	GROUP BY match_id;
	`

	row := repo.db.QueryRowContext(ctx, query, matchID)

	output := entity.MatchStats{}

	if err := row.Scan(
		&output.MatchID,
		&output.AvgRate,
		&output.Rate1,
		&output.Rate2,
		&output.Rate3,
		&output.Rate4,
		&output.Rate5,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &output, nil
}
