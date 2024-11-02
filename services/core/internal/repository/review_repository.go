package repository

import (
	"database/sql"
	"fmt"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/utils"
)

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (repo *reviewRepository) Create(review *domain.Review) error {
	query := `
		INSERT INTO reviews (user_id, rate, description, match_id, home_team_id, away_team_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := repo.db.Exec(
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

func (repo *reviewRepository) Delete(reviewID int64) error {
	query := `
		DELETE FROM reviews WHERE id = $1
	`

	if _, err := repo.db.Exec(query, reviewID); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) FindOneByID(requesterID string, reviewID int64) (*domain.Review, error) {
	query := `
		SELECT 
			r.id as reviewID, r.rate, r.description, 
			r.match_id, r.created_at, u.id, u.username, u.favorite_team,
			COUNT(l.review_id) as like_count,
			COUNT(CASE WHEN l.like_owner_id = $1 THEN 1 END) > 0 AS is_liked
		FROM reviews r
		LEFT JOIN users u ON u.id = r.user_id
		LEFT JOIN likes l ON l.review_id = reviewID
		WHERE r.id = $2
	`

	row := repo.db.QueryRow(query, requesterID, reviewID)

	review, err := repo.scanReview(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr(
				"review não encontrada",
				404,
				"REVIEW_REPOSITORY:FIND_REVIEW_BY_ID:REVIEW_NOT_FOUND",
			)
		}
		return nil, errors.NewHTTPErr(
			"ocorreu um erro desconhecido",
			500,
			"REVIEW_REPOSITORY:FIND_REVIEW_BY_ID:SQL_ERROR",
		)
	}

	return review, nil
}

func (repo *reviewRepository) ExistsByID(reviewID int64) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM reviews WHERE reviews.id = $1)
	`

	row := repo.db.QueryRow(query, reviewID)

	var exists *sql.NullBool

	if err := row.Scan(&exists); err != nil {

		return false, err
	}

	return exists.Bool, nil
}

func (repo *reviewRepository) ListFeed(requesterID, strategy string, page *pagination.Page) ([]domain.Review, error) {
	orderBy := func() string {
		switch strategy {
		case "relevant":
			return ""
		case "newer":
			return "ORDER BY created_at DESC"
		case "older":
			return "ORDER BY created_at"
		default:
			return ""
		}
	}()

	query := fmt.Sprintf(`
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
				r.away_team_id
			FROM reviews r
			WHERE r.created_at >= NOW() - INTERVAL '7 days'
			LIMIT 500
		),
		ranked_feed AS (
			SELECT 
				rr.reviewID, 
				rr.rate, 
				rr.description, 
				rr.match_id, 
				rr.created_at,
				u.id,
				u.username, 
				u.favorite_team,
				COUNT(l.review_id) as like_count,
      			COUNT(CASE WHEN l.like_owner_id = $1 THEN 1 END) > 0 AS is_liked
			FROM recent_reviews rr
			JOIN users u ON u.id = rr.user_id
			LEFT JOIN requester req ON true
			LEFT JOIN followers f ON f.following_id = rr.user_id AND f.follower_id = req.id
			LEFT JOIN likes l ON l.review_id = rr.reviewID
			GROUP BY 
				rr.reviewID, rr.rate, rr.description, rr.match_id, rr.created_at, 
				u.id, u.username, u.favorite_team, f.follower_id, req.favorite_team,
				rr.home_team_id, rr.away_team_id
			ORDER BY (CASE 
					WHEN f.follower_id IS NOT NULL THEN 1.5
					ELSE 0 
				END) +
				(CASE
					WHEN rr.home_team_id = req.favorite_team THEN 0.6
					WHEN rr.away_team_id = req.favorite_team THEN 0.5
					ELSE 0
				END) + 
				(CASE
					WHEN rr.created_at >= NOW() - INTERVAL '1 day' THEN 0.30
					WHEN rr.created_at >= NOW() - INTERVAL '2 days' THEN 0.20
					WHEN rr.created_at >= NOW() - INTERVAL '3 days' THEN 0.10
					ELSE 0
				END) DESC
		)

		SELECT * FROM ranked_feed
		%s
		LIMIT $2
		OFFSET ($3 - 1) * $2
	`, orderBy)

	rows, err := repo.db.Query(query, requesterID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	reviews, err := repo.scanReviews(rows)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (repo *reviewRepository) ListAll(requesterID, userID, team, match string, page *pagination.Page) ([]domain.Review, error) {

	whereBuilder := utils.NewWhereBuilder()

	whereBuilder.AddFilter("user_id", userID, "AND")
	whereBuilder.AddFilter("match_id", match, "AND")

	whereBuilder.AddGroupFilter(
		"AND",
		whereBuilder.NewFilter("home_team_id", team, "OR"),
		whereBuilder.NewFilter("away_team_id", team, "OR"),
	)

	whereQuery, params := whereBuilder.Build()

	query := fmt.Sprintf(`
		SELECT 
			r.id as reviewID, 
			r.rate, 
			r.description, 
			r.match_id, 
			r.created_at,
			u.id, 
			u.username, 
			u.favorite_team,
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
	`, len(params)+1, whereQuery, len(params)+2, len(params)+3, len(params)+2)

	rows, err := repo.db.Query(query, append(params, requesterID, page.Size, page.Index)...)
	if err != nil {
		return nil, err
	}

	return repo.scanReviews(rows)
}

func (repo *reviewRepository) Like(requesterID string, reviewID int64) error {
	query := `
		INSERT INTO likes (like_owner_id, review_id)
		VALUES ($1, $2)
	`

	if _, err := repo.db.Exec(query, requesterID, reviewID); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) Unlike(requesterID string, reviewID int64) error {
	query := `
		DELETE FROM likes WHERE like_owner_id = $1 AND review_id = $2
	`

	if _, err := repo.db.Exec(query, requesterID, reviewID); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) IsLiked(requesterID string, reviewID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM likes
			WHERE like_owner_id = $1 AND review_id = $2
		)
	`

	var exists bool
	if err := repo.db.QueryRow(query, requesterID, reviewID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *reviewRepository) scanReview(row *sql.Row) (*domain.Review, error) {
	var review domain.Review

	review.Author = &domain.User{}

	if err := row.Scan(
		&review.ID, &review.Rate, &review.Description,
		&review.MatchID, &review.CreatedAt, &review.Author.ID,
		&review.Author.Username,
		&review.Author.FavoriteTeamID,
		&review.Likes,
		&review.IsLiked,
	); err != nil {
		return nil, err
	}

	return &review, nil
}

func (repo *reviewRepository) scanReviews(rows *sql.Rows) ([]domain.Review, error) {
	defer rows.Close()

	var output []domain.Review

	for rows.Next() {
		var review domain.Review

		review.Author = &domain.User{}

		if err := rows.Scan(
			&review.ID, &review.Rate, &review.Description,
			&review.MatchID, &review.CreatedAt, &review.Author.ID,
			&review.Author.Username,
			&review.Author.FavoriteTeamID,
			&review.Likes,
			&review.IsLiked,
		); err != nil {
			return nil, err
		}

		output = append(output, review)
	}

	return output, nil
}
