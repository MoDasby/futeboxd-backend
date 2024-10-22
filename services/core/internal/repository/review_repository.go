package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/pkg/utils"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (repo *reviewRepository) AddReview(review *domain.Review) error {
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

func (repo *reviewRepository) DeleteReview(reviewID int64) error {
	query := `
		DELETE FROM reviews WHERE id = $1
	`

	if _, err := repo.db.Exec(query, reviewID); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) FindReviewByID(reviewID int64) (*domain.Review, error) {
	query := `
		SELECT r.id as reviewID, r.rate, r.description, r.match_id, r.created_at, u.id, u.username, u.favorite_team
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		WHERE r.id = $1
	`

	row := repo.db.QueryRow(query, reviewID)

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

func (repo *reviewRepository) ListFeed(requesterID, strategy string, pageSize, pageIndex int) ([]domain.Review, error) {
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
			ORDER BY r.created_at DESC
			LIMIT 1000
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
				u.favorite_team
			FROM recent_reviews rr
			JOIN users u ON u.id = rr.user_id
			CROSS JOIN requester req
			LEFT JOIN followers f ON f.following_id = rr.user_id AND f.follower_id = req.id
			ORDER BY (CASE 
					WHEN f.follower_id IS NOT NULL THEN 1.5
					ELSE 0 
				END) +
				(CASE
					WHEN rr.home_team_id = req.favorite_team THEN 0.6
					WHEN rr.away_team_id = req.favorite_team THEN 0.5
					ELSE 0
				END) DESC, rr.created_at DESC
		)

		SELECT * FROM ranked_feed
		%s
		LIMIT $2
		OFFSET ($3 - 1) * $2
	`, orderBy)

	rows, err := repo.db.Query(query, requesterID, pageSize, pageIndex)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	reviews, err := repo.scanReviews(rows)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return reviews, nil
}

func (repo *reviewRepository) ListReviews(pageSize, pageIndex int, userID, team, match string) ([]domain.Review, error) {

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
			u.favorite_team
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		%s
		ORDER BY r.created_at DESC
		LIMIT $%d
		OFFSET ($%d - 1) * $%d
	`, whereQuery, len(params)+1, len(params)+2, len(params)+1)

	rows, err := repo.db.Query(query, append(params, pageSize, pageIndex)...)
	if err != nil {
		return nil, err
	}

	return repo.scanReviews(rows)
}

func (repo *reviewRepository) scanReview(row *sql.Row) (*domain.Review, error) {
	var review domain.Review

	review.Author = &domain.User{}

	if err := row.Scan(
		&review.ID, &review.Rate, &review.Description,
		&review.MatchID, &review.CreatedAt, &review.Author.ID,
		&review.Author.Username,
		&review.Author.FavoriteTeamID,
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
		); err != nil {
			return nil, err
		}

		output = append(output, review)
	}

	return output, nil
}
