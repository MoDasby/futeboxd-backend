package review

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/review/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type ListReviewsOptions struct {
	Username   string
	Team       string
	Match      string
	SearchTerm string
}

type ReviewUsecases interface {
	Create(ctx context.Context, input *dto.ReviewInput) error
	Delete(ctx context.Context, reviewID int64) error
	ListFeed(ctx context.Context, page *pagination.Page) ([]dto.Review, error)
	ListBy(ctx context.Context, options ListReviewsOptions, page *pagination.Page) ([]dto.Review, error)
	Search(ctx context.Context, term string, page *pagination.Page) ([]dto.Review, error)
	ToggleLike(ctx context.Context, reviewID int64) (*dto.LikeStats, error)
}
