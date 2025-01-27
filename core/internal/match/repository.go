package match

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/match/entity"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type Repository interface {
	ListPopularMatches(ctx context.Context, page *pagination.Page) ([]int64, error)
	GetMatchStats(ctx context.Context, matchID int64) (*entity.MatchStats, error)
}
