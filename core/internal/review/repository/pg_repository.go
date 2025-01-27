package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/modasby/futeboxd-backend/core/internal/review"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) review.Repository {
	return &reviewRepository{db: db}
}

func (repo *reviewRepository) Upsert(ctx context.Context, review *review.Review) error {
	query := `
		INSERT INTO reviews (user_id, rate, description, match_id, home_team_id, away_team_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, match_id) 
		DO UPDATE SET 
			rate = EXCLUDED.rate, 
			description = EXCLUDED.description,
			home_team_id = EXCLUDED.home_team_id,
			away_team_id = EXCLUDED.away_team_id
	`

	_, err := repo.db.ExecContext(
		ctx,
		query,
		review.Author.ID,
		review.Rate,
		review.Description,
		review.MatchID,
		review.HomeTeamID,
		review.AwayTeamID,
	)

	return err
}

func (repo *reviewRepository) Delete(ctx context.Context, reviewID int64) error {
	query := `
		DELETE FROM reviews WHERE id = $1
	`

	if _, err := repo.db.ExecContext(ctx, query, reviewID); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) FindOneByID(ctx context.Context, requesterID string, reviewID int64) (*review.Review, error) {
	query := `
		SELECT 
			r.id, r.rate, r.description, 
			r.match_id, r.created_at, u.id as user_id, u.name, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM comments c WHERE c.review_id = r.id) as comments_count,
			COUNT(l.review_id) as like_count,
			COUNT(CASE WHEN l.like_owner_id = $1 THEN 1 END) > 0 AS is_liked
		FROM reviews r
		LEFT JOIN users u ON u.id = r.user_id
		LEFT JOIN likes l ON l.review_id = r.id
		WHERE r.id = $2
		GROUP BY r.id, u.id
	`

	row := repo.db.QueryRowContext(ctx, query, requesterID, reviewID)

	review, err := repo.scanReview(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr(
				"review não encontrada",
				404,
				"REVIEW_REPOSITORY:FIND_REVIEW_BY_ID:REVIEW_NOT_FOUND",
			)
		}
		return nil, err
	}

	return review, nil
}

func (repo *reviewRepository) ExistsByID(ctx context.Context, reviewID int64) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM reviews WHERE reviews.id = $1)
	`

	row := repo.db.QueryRowContext(ctx, query, reviewID)

	var exists *sql.NullBool

	if err := row.Scan(&exists); err != nil {

		return false, err
	}

	return exists.Bool, nil
}

func (repo *reviewRepository) ListFeed(ctx context.Context, requesterID string, page *pagination.Page) ([]review.Review, error) {

	query := `
		WITH requester AS (
			SELECT id, favorite_team FROM users WHERE id = $1
		),
		recent_reviews AS (
			SELECT 
				r.id as reviewID, 
				r.rate, 
				r.description, 
				r.match_id, 
				r.created_at,
				r.user_id,
				r.home_team_id,
				r.away_team_id,
				(SELECT COUNT(*) FROM comments c WHERE c.review_id = r.id) as comments_count,
				COUNT(l.review_id) as like_count,
				COUNT(CASE WHEN l.like_owner_id = $1 THEN 1 END) > 0 AS is_liked
			FROM reviews r
			INNER JOIN followers f ON f.following_id = r.user_id
			LEFT JOIN likes l ON l.review_id = r.id
			CROSS JOIN requester req
			WHERE f.follower_id = req.id
			GROUP BY
				r.id
			LIMIT 500
		)
		SELECT 
			rr.reviewID, 
			rr.rate, 
			rr.description, 
			rr.match_id, 
			rr.created_at,
			u.id as user_id, u.name, u.username, u.favorite_team, u.profile_picture,
			rr.comments_count,
			rr.like_count,
			rr.is_liked
		FROM recent_reviews rr
		JOIN users u ON u.id = rr.user_id
		LEFT JOIN requester req ON true
		GROUP BY 
			rr.reviewID, rr.rate, rr.description, rr.match_id, rr.created_at, 
			u.id, u.username, u.favorite_team, req.favorite_team,
			rr.home_team_id, rr.away_team_id, rr.comments_count, rr.like_count, rr.is_liked
		ORDER BY
			(CASE
				WHEN rr.home_team_id = req.favorite_team THEN 1
				WHEN rr.away_team_id = req.favorite_team THEN 0.7
				ELSE 0
			END) + (CASE
				WHEN rr.comments_count > 0 THEN 0.2
				WHEN rr.comments_count > 10 THEN 0.4
				WHEN rr.comments_count > 100 THEN 0.5
				WHEN rr.comments_count > 1000 THEN 0.8
				WHEN rr.comments_count > 10000 THEN 1
			END) + (CASE
				WHEN rr.like_count > 0 THEN 0.2
				WHEN rr.like_count > 10 THEN 0.4
				WHEN rr.like_count > 100 THEN 0.5
				WHEN rr.like_count > 1000 THEN 0.8
				WHEN rr.like_count > 10000 THEN 1
			END), rr.created_at DESC
		LIMIT $2
		OFFSET ($3 - 1) * $2
	`

	rows, err := repo.db.QueryContext(ctx, query, requesterID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	reviews, err := repo.scanReviews(rows)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (repo *reviewRepository) ListAll(ctx context.Context, requesterID, where string, params []any, page *pagination.Page) ([]review.Review, error) {

	query := fmt.Sprintf(`
		SELECT
			r.id as reviewID, 
			r.rate, 
			r.description, 
			r.match_id, 
			r.created_at,
			u.id as user_id, u.name, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM comments c WHERE c.review_id = r.id) as comments_count,
			COUNT(l.review_id) as like_count,
			COUNT(CASE WHEN l.like_owner_id = $%d THEN 1 END) > 0 AS is_liked
		FROM reviews r
		LEFT JOIN users u ON u.id = r.user_id
		LEFT JOIN likes l ON l.review_id = r.id
		%s
		GROUP BY r.id, u.id
		ORDER BY r.created_at DESC
		LIMIT $%d
		OFFSET ($%d - 1) * $%d
	`, len(params)+1, where, len(params)+2, len(params)+3, len(params)+2)

	rows, err := repo.db.QueryContext(ctx, query, append(params, requesterID, page.Size, page.Index)...)
	if err != nil {
		return nil, err
	}

	return repo.scanReviews(rows)
}

func (repo *reviewRepository) Search(ctx context.Context, requesterID, term string, page *pagination.Page) ([]review.Review, error) {
	query := `
		SELECT 
			r.id as reviewID, 
			r.rate, 
			r.description, 
			r.match_id, 
			r.created_at,
			u.id as user_id, u.name, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM comments c WHERE c.review_id = r.id) as comments_count,
			COUNT(l.review_id) as like_count,
			COUNT(CASE WHEN l.like_owner_id = $2 THEN 1 END) > 0 AS is_liked
		FROM reviews r
		LEFT JOIN users u ON u.id = r.user_id
		LEFT JOIN likes l ON l.review_id = r.id
		WHERE r.search_vector @@ websearch_to_tsquery('portuguese', $1)
		GROUP BY r.id, u.id
		LIMIT $3
		OFFSET ($4 - 1) * $3
	`

	rows, err := repo.db.QueryContext(ctx, query, term, requesterID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews, err := repo.scanReviews(rows)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (repo *reviewRepository) Like(ctx context.Context, requesterID string, reviewID int64) error {
	query := `
		INSERT INTO likes (like_owner_id, review_id)
		VALUES ($1, $2)
	`

	if _, err := repo.db.ExecContext(ctx, query, requesterID, reviewID); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) Unlike(ctx context.Context, requesterID string, reviewID int64) error {
	query := `
		DELETE FROM likes WHERE like_owner_id = $1 AND review_id = $2
	`

	if _, err := repo.db.ExecContext(ctx, query, requesterID, reviewID); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) IsLiked(ctx context.Context, requesterID string, reviewID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM likes
			WHERE like_owner_id = $1 AND review_id = $2
		)
	`

	var exists bool
	if err := repo.db.QueryRowContext(ctx, query, requesterID, reviewID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *reviewRepository) scanReview(row *sql.Row) (*review.Review, error) {
	var review review.Review
	var description sql.NullString

	review.Author = &user.User{}

	if err := row.Scan(
		&review.ID, &review.Rate, &description,
		&review.MatchID, &review.CreatedAt,
		&review.Author.ID, &review.Author.Name, &review.Author.Username,
		&review.Author.FavoriteTeamID, &review.Author.ProfilePicture,
		&review.CommentsCount,
		&review.Likes,
		&review.IsLiked,
	); err != nil {
		return nil, err
	}

	if description.Valid {
		review.Description = description.String
	}

	return &review, nil
}

func (repo *reviewRepository) scanReviews(rows *sql.Rows) ([]review.Review, error) {
	defer rows.Close()

	var output []review.Review

	for rows.Next() {
		var review review.Review
		var description sql.NullString

		review.Author = &user.User{}

		if err := rows.Scan(
			&review.ID, &review.Rate, &description,
			&review.MatchID, &review.CreatedAt,
			&review.Author.ID, &review.Author.Name, &review.Author.Username,
			&review.Author.FavoriteTeamID, &review.Author.ProfilePicture,
			&review.CommentsCount,
			&review.Likes,
			&review.IsLiked,
		); err != nil {
			return nil, err
		}

		if description.Valid {
			review.Description = description.String
		}

		output = append(output, review)
	}

	return output, nil
}
