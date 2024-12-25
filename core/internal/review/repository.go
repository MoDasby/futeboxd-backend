package review

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type Repository interface {
	Create(ctx context.Context, review *Review) error
	Delete(ctx context.Context, reviewID int64) error
	ExistsByID(ctx context.Context, reviewID int64) (bool, error)
	ListFeed(ctx context.Context, requesterID string, page *pagination.Page) ([]Review, error)
	ListAll(ctx context.Context, requesterID, where string, params []any, page *pagination.Page) ([]Review, error)
	ListTrendingMatches(ctx context.Context, page *pagination.Page) ([]int64, error)
	Search(ctx context.Context, requesterID, term string, page *pagination.Page) ([]Review, error)
	FindOneByID(ctx context.Context, requesterID string, reviewID int64) (*Review, error)
	Like(ctx context.Context, requesterID string, reviewID int64) error
	Unlike(ctx context.Context, requesterID string, reviewID int64) error
	IsLiked(ctx context.Context, requesterID string, reviewID int64) (bool, error)
}
