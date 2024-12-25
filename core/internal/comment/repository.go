package comment

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, commentID int64) error
	ExistsByID(ctx context.Context, commentID int64) (bool, error)
	FindOneByID(ctx context.Context, requesterID string, commentID int64) (*Comment, error)
	ListByReview(ctx context.Context, requesterID string, reviewID int64, page *pagination.Page) ([]Comment, error)
	Like(ctx context.Context, requesterID string, commentID int64) error
	Unlike(ctx context.Context, requesterID string, commentID int64) error
}
