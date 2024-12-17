package match

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/domain"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type MatchUsecases interface {
	ListPopularMatches(ctx context.Context, page *pagination.Page) ([]football.Match, error)
}

type matchUsecases struct {
	matchRepo      domain.MatchRepository
	footballClient football.Client
}

func NewMatchUsecases(
	matchRepo domain.MatchRepository,
	footballClient football.Client,
) MatchUsecases {
	return &matchUsecases{
		matchRepo:      matchRepo,
		footballClient: footballClient,
	}
}

func (uc *matchUsecases) ListPopularMatches(ctx context.Context, page *pagination.Page) ([]football.Match, error) {
	popularMatchesIds, err := uc.matchRepo.ListPopularMatches(ctx, page)
	if err != nil {
		return nil, err
	}

	return uc.footballClient.GetMatches(popularMatchesIds)
}
