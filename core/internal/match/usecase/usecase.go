package usecase

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/match"
	"github.com/modasby/futeboxd-backend/core/internal/match/dto"
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

func (uc *matchUsecases) GetMatchStats(ctx context.Context, matchID int64) (*dto.MatchStats, error) {
	match, err := uc.footballClient.GetMatch(matchID)
	if err != nil {
		return nil, err
	}

	matchStats, err := uc.matchRepo.GetMatchStats(ctx, matchID)
	if err != nil {
		return nil, err
	}

	if matchStats == nil {
		return &dto.MatchStats{
			Match:   match,
			AvgRate: 0,
			ReviewsSummary: dto.ReviewsSummary{
				Count1: 0,
				Count2: 0,
				Count3: 0,
				Count4: 0,
				Count5: 0,
			},
		}, nil
	}

	return &dto.MatchStats{
		Match:   match,
		AvgRate: matchStats.AvgRate,
		ReviewsSummary: dto.ReviewsSummary{
			Count1: matchStats.Rate1,
			Count2: matchStats.Rate2,
			Count3: matchStats.Rate3,
			Count4: matchStats.Rate4,
			Count5: matchStats.Rate5,
		},
	}, nil
}
