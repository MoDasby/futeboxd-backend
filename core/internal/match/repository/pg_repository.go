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

func (repo *matchRepository) ListPopularMatches(ctx context.Context, page *pagination.Page) ([]entity.MatchStats, error) {
	query := `
		WITH global_avg AS (
			SELECT avg(rate) AS m from reviews
		),
		match_reviews AS (
			SELECT
				r.match_id,
				COUNT(r.match_id) AS n,
				SUM(r.rate) AS sum_r,
				(SELECT m FROM global_avg) AS m,
				COUNT(CASE WHEN rate = 1 THEN 1 END) as count1,
				COUNT(CASE WHEN rate = 2 THEN 1 END) as count2,
				COUNT(CASE WHEN rate = 3 THEN 1 END) as count3,
				COUNT(CASE WHEN rate = 4 THEN 1 END) as count4,
				COUNT(CASE WHEN rate = 5 THEN 1 END) as count5
			FROM reviews r
			WHERE r.created_at >= NOW() - INTERVAL '30 days'
			GROUP BY r.match_id
		)

		SELECT
			match_id,
			ROUND((20 * m + sum_r) / (20 + n), 2) as avg_rate,
			count1,
			count2,
			count3,
			count4,
			count5
		FROM match_reviews
		ORDER BY avg_rate DESC
		LIMIT $1
		OFFSET ($2 - 1) * $1
	`

	var matches []entity.MatchStats

	rows, err := repo.db.QueryContext(ctx, query, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var match entity.MatchStats
		if err := rows.Scan(
			&match.MatchID,
			&match.AvgRate,
			&match.Rate1,
			&match.Rate2,
			&match.Rate3,
			&match.Rate4,
			&match.Rate5,
		); err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (repo *matchRepository) GetMatchStats(ctx context.Context, matchID int64) (*entity.MatchStats, error) {
	query := `
	SELECT 
		match_id,
		ROUND(AVG(rate), 2) as avg_rate,
		COUNT(CASE WHEN rate = 1 THEN 1 END) as count1,
		COUNT(CASE WHEN rate = 2 THEN 1 END) as count2,
		COUNT(CASE WHEN rate = 3 THEN 1 END) as count3,
		COUNT(CASE WHEN rate = 4 THEN 1 END) as count4,
		COUNT(CASE WHEN rate = 5 THEN 1 END) as count5
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
