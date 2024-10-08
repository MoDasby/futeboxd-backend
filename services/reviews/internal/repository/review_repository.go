package repository

import (
	"database/sql"
	"fmt"

	"github.com/modasby/futeboxd-api/pkg/utils"
	"github.com/modasby/futeboxd-api/services/reviews/internal/domain"
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
		review.UserID,
		review.Rate,
		review.Description,
		review.MatchID,
		review.HomeTeamID,
		review.AwayTeamID,
	)

	return err
}

func (repo *reviewRepository) ListReviews(pageSize, pageIndex int, userID, team, match string) (*utils.Pageable[*domain.Review], error) {

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
		WITH total_count AS (
			SELECT COUNT(*) AS total_reviews FROM reviews %s
		)

		SELECT r.id as reviewID, r.rate, r.description, r.match_id, r.user_id, 
		total_count.total_reviews as items_count
		FROM reviews r
		JOIN total_count ON TRUE
		%s
		ORDER BY r.created_at DESC
		LIMIT $%d
		OFFSET ($%d - 1) * $%d
	`, whereQuery, whereQuery, len(params)+1, len(params)+2, len(params)+1)

	fmt.Println(query)

	rows, error := repo.db.Query(query, append(params, pageSize, pageIndex)...)
	if error != nil {
		return nil, error
	}
	defer rows.Close()

	var output utils.Pageable[*domain.Review]

	output.Items = make([]*domain.Review, 0)

	var itemsCount sql.NullInt32

	for rows.Next() {
		var description, userID sql.NullString
		var reviewID, rate sql.NullInt32
		var matchID sql.NullInt64

		if err := rows.Scan(
			&reviewID, &rate, &description,
			&matchID, &userID, &itemsCount,
		); err != nil {
			return nil, error
		}

		review := domain.Review{
			ID:          int(reviewID.Int32),
			UserID:      userID.String,
			Description: description.String,
			Rate:        int(rate.Int32),
			MatchID:     matchID.Int64,
		}

		output.Items = append(output.Items, &review)
	}

	output.CurrentPage = pageIndex
	output.Size = pageSize

	return &output, nil
}
