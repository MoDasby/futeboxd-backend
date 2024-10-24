package repository

import (
	"database/sql"

	"github.com/modasby/futeboxd-api/pkg/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type CommentsRepository struct {
	db *sql.DB
}

func NewCommentsRepository(
	db *sql.DB,
) domain.CommentsRepository {
	return &CommentsRepository{
		db: db,
	}
}

func (repo *CommentsRepository) Create(comment *domain.Comment) error {
	query := `
		INSERT INTO comments (user_id, review_id, content)
		VALUES ($1, $2, $3)
	`

	if _, err := repo.db.Exec(query, comment.Author.ID, comment.ParentID, comment.Content); err != nil {
		return err
	}

	return nil
}
func (repo *CommentsRepository) Delete(commentID int64) error {
	query := `
		DELETE FROM comments WHERE id = $1
	`

	if _, err := repo.db.Exec(query, commentID); err != nil {
		return err
	}

	return nil
}
func (repo *CommentsRepository) ListByReview(reviewID int64, page *pagination.Page) ([]domain.Comment, error) {
	query := `
		SELECT 
			u.id as user_id, u.username, u.email, u.favorite_team,
			c.id as comment_id, c.content, c.created_at
		FROM comments c
		LEFT JOIN users u ON u.id = c.user_id
		WHERE c.review_id = $1
		LIMIT $2
		OFFSET ($3 - 1) * $2
	`

	rows, err := repo.db.Query(query, reviewID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	output := make([]domain.Comment, 0)

	for rows.Next() {
		var comment domain.Comment

		if err := rows.Scan(
			&comment.Author.ID, &comment.Author.Username, &comment.Author.Email,
			&comment.Author.FavoriteTeamID, &comment.ID, &comment.Content, &comment.CreatedAt,
		); err != nil {
			return nil, err
		}

		output = append(output, comment)
	}

	return output, nil
}
