package repository

import (
	"context"
	"database/sql"

	"github.com/modasby/futeboxd-backend/core/internal/domain"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
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

func (repo *CommentsRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (user_id, review_id, content)
		VALUES ($1, $2, $3)
	`

	if _, err := repo.db.ExecContext(ctx, query, comment.Author.ID, comment.ParentID, comment.Content); err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepository) Delete(ctx context.Context, commentID int64) error {
	query := `
		DELETE FROM comments WHERE id = $1
	`

	if _, err := repo.db.ExecContext(ctx, query, commentID); err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepository) ExistsByID(ctx context.Context, commentID int64) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM comments WHERE comments.id = $1)
	`

	row := repo.db.QueryRowContext(ctx, query, commentID)

	var exists *sql.NullBool

	if err := row.Scan(&exists); err != nil {

		return false, err
	}

	return exists.Bool, nil
}

func (repo *CommentsRepository) FindOneByID(ctx context.Context, requesterID string, commentID int64) (*domain.Comment, error) {
	query := `
		SELECT 
			c.id as comment_id, u.id as user_id, 
			c.review_id, c.content, c.created_at,
			COUNT(l.comment_id) as like_count,
      		COUNT(CASE WHEN l.like_owner_id = $2 THEN 1 END) > 0 AS is_liked
		FROM comments c
		LEFT JOIN users u ON u.id = c.user_id
		LEFT JOIN likes l ON l.comment_id = c.id
		WHERE c.id = $1
		GROUP BY 
			u.id, u.username, u.email, u.favorite_team,
			c.id, c.content, c.created_at
	`

	row := repo.db.QueryRowContext(ctx, query, commentID, requesterID)

	var comment domain.Comment

	var author domain.User

	comment.Author = &author

	if err := row.Scan(
		&comment.ID, &comment.Author.ID, &comment.ParentID,
		&comment.Content, &comment.CreatedAt, &comment.LikeCount,
		&comment.IsLiked,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr(
				"comentário não encontrado",
				404,
				"COMMENT_REPOSITORY:FIND_COMMENT_BY_ID:COMMENT_NOT_FOUND",
			)
		}
		return nil, err
	}

	return &comment, nil
}

func (repo *CommentsRepository) ListByReview(
	ctx context.Context, requesterID string,
	reviewID int64,
	page *pagination.Page,
) ([]domain.Comment, error) {
	query := `
		SELECT 
			u.id as user_id, u.username, u.email, u.favorite_team,
			c.id as comment_id, c.content, c.created_at,
			COUNT(l.comment_id) as like_count,
      		COUNT(CASE WHEN l.like_owner_id = $1 THEN 1 END) > 0 AS is_liked
		FROM comments c
		LEFT JOIN users u ON u.id = c.user_id
		LEFT JOIN likes l ON l.comment_id = c.id
		WHERE c.review_id = $2
		GROUP BY 
			u.id, u.username, u.email, u.favorite_team,
			c.id, c.content, c.created_at
		ORDER BY c.created_at DESC
		LIMIT $3
		OFFSET ($4 - 1) * $3
	`

	rows, err := repo.db.QueryContext(ctx, query, requesterID, reviewID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	output := make([]domain.Comment, 0)

	for rows.Next() {
		var comment domain.Comment
		var author domain.User

		comment.Author = &author

		if err := rows.Scan(
			&comment.Author.ID, &comment.Author.Username, &comment.Author.Email,
			&comment.Author.FavoriteTeamID, &comment.ID, &comment.Content, &comment.CreatedAt,
			&comment.LikeCount, &comment.IsLiked,
		); err != nil {
			return nil, err
		}

		output = append(output, comment)
	}

	return output, nil
}

func (repo *CommentsRepository) Like(ctx context.Context, requesterID string, commentID int64) error {
	query := `
		INSERT INTO likes (like_owner_id, comment_id)
		VALUES ($1, $2)
	`

	if _, err := repo.db.ExecContext(ctx, query, requesterID, commentID); err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepository) Unlike(ctx context.Context, requesterID string, commentID int64) error {
	query := `
		DELETE FROM likes WHERE like_owner_id = $1 AND comment_id = $2
	`

	if _, err := repo.db.ExecContext(ctx, query, requesterID, commentID); err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepository) IsLiked(ctx context.Context, requesterID string, commentID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM likes
			WHERE like_owner_id = $1 AND comment_id = $2
		)
	`

	var exists bool
	if err := repo.db.QueryRowContext(ctx, query, requesterID, commentID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
