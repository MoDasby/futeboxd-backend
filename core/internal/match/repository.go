package match

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type Repository interface {
	ListPopularMatches(ctx context.Context, page *pagination.Page) ([]int64, error)
}
