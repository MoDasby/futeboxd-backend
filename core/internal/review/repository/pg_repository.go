package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/review"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
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
			return nil, &errors.HTTPErr{
				Msg:        "review não encontrada",
				Code:       http.StatusNotFound,
				Context:    "REVIEW:REPOSITORY:FIND_REVIEW_BY_ID:REVIEW_NOT_FOUND",
				StackTrace: errors.CaptureStackTrace(),
				ErrorCode:  utils.GetTraceIDFromCtx(ctx),
				Timestamp:  time.Now().UTC(),
			}
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
		candidate_reviews AS (
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
				(SELECT COUNT(*) FROM likes l WHERE l.review_id = r.id) as like_count,
				EXISTS((SELECT 1 FROM likes l WHERE l.review_id = r.id AND l.like_owner_id = $1)) AS is_liked,
				EXISTS(SELECT 1 FROM followers f WHERE f.follower_id = $1 AND f.following_id = r.user_id) AS is_following,
				GREATEST(EXTRACT(EPOCH FROM NOW() - r.created_at) / 3600, 1) AS age_hours
			FROM reviews r
			JOIN requester req ON true
			WHERE 
				r.user_id != $1
				AND (
					EXISTS (
						SELECT 1 FROM followers f WHERE f.follower_id = $1 AND f.following_id = r.user_id
					)
					OR req.favorite_team IN (r.home_team_id, r.away_team_id)
					OR (
						(SELECT COUNT(*) FROM likes l WHERE l.review_id = r.id) >= 5
						OR (SELECT COUNT(*) FROM comments c WHERE c.review_id = r.id) >= 3
					)
				)
			ORDER BY r.created_at DESC
			LIMIT 500
		)
		SELECT 
			cr.reviewID, 
			cr.rate, 
			cr.description, 
			cr.match_id, 
			cr.created_at,
			u.id as user_id, u.name, u.username, u.favorite_team, u.profile_picture,
			cr.comments_count,
			cr.like_count,
			cr.is_liked
		FROM candidate_reviews cr
		JOIN users u ON u.id = cr.user_id
		JOIN requester req ON true
		ORDER BY
			(
				-- Engagement velocity: likes+comments per hour (Twitter-like trending signal)
				LOG(2, GREATEST(cr.like_count + cr.comments_count * 2, 1) + 1)
				/ POWER(cr.age_hours, 0.3)
			)

			-- Team relevance boost
			* (CASE
				WHEN cr.home_team_id = req.favorite_team THEN 1.3
				WHEN cr.away_team_id = req.favorite_team THEN 1.2
				ELSE 1
			END)

			-- Social graph boost (following)
			* (CASE 
				WHEN cr.is_following THEN 1.5
				ELSE 0.8
			END)

			-- Engagement tier bonus (fixed ordering: highest thresholds first)
			* (CASE
				WHEN cr.like_count >= 100 THEN 1.8
				WHEN cr.like_count >= 50 THEN 1.5
				WHEN cr.like_count >= 10 THEN 1.3
				WHEN cr.like_count >= 3 THEN 1.1
				ELSE 1
			END)

			-- Comment engagement tier bonus
			* (CASE
				WHEN cr.comments_count >= 50 THEN 1.6
				WHEN cr.comments_count >= 20 THEN 1.4
				WHEN cr.comments_count >= 5 THEN 1.2
				WHEN cr.comments_count >= 1 THEN 1.1
				ELSE 1
			END)

			-- Content quality: reviews with descriptions are more valuable
			* (CASE
				WHEN cr.description IS NOT NULL AND LENGTH(cr.description) > 100 THEN 1.3
				WHEN cr.description IS NOT NULL AND LENGTH(cr.description) > 0 THEN 1.1
				ELSE 1
			END)

			-- Gentle time decay (much less aggressive than before)
			* POWER(0.95, EXTRACT(EPOCH FROM NOW() - cr.created_at) / 86400)

			-- Small random factor for content discovery
			+ RANDOM() * 0.1

			DESC
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
