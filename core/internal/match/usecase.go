package match

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/match/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type MatchUsecases interface {
	ListPopularMatches(ctx context.Context, page *pagination.Page) ([]football.Match, error)
	GetMatchStats(ctx context.Context, matchID int64) (*dto.MatchStats, error)
}
