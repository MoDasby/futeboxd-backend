package usecase

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/match"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type matchUsecases struct {
	matchRepo      match.Repository
	footballClient football.Client
}

func NewMatchUsecases(
	matchRepo match.Repository,
	footballClient football.Client,
) match.MatchUsecases {
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
