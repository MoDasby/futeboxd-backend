package comment

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/comment/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type CommentUsecases interface {
	Create(ctx context.Context, input *dto.CommentInput) error
	Delete(ctx context.Context, commentID int64) error
	ListByReview(ctx context.Context, reviewID int64, page *pagination.Page) ([]dto.Comment, error)
	ToggleLike(ctx context.Context, commentID int64) (*dto.LikeStats, error)
}
