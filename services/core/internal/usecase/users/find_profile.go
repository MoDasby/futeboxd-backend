package usecase

import (
	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/queries"
)

type FindProfile struct {
	userRepository          domain.UserRepository
	footballClient          football.Client
	followStatsQueryService queries.FollowStatsQuery
}

func NewFindProfile(
	userRepository domain.UserRepository,
	footballClient football.Client,
	followStatsQueryService queries.FollowStatsQuery,
) *FindProfile {
	return &FindProfile{
		userRepository:          userRepository,
		footballClient:          footballClient,
		followStatsQueryService: followStatsQueryService,
	}
}

func (uc *FindProfile) Execute(requester *domain.User, username string) (*dto.ProfileDTO, error) {
	following, err := uc.userRepository.FindOneByIdOrUsername(username)
	if err != nil {
		return nil, err
	}

	followStats, err := uc.followStatsQueryService.GetFollowStats(following.ID, requester.ID)
	if err != nil {
		return nil, err
	}

	favoriteTeam, err := uc.footballClient.GetTeam(following.FavoriteTeamID)
	if err != nil {
		return nil, err
	}

	output := dto.ProfileDTO{
		UserDTO: &dto.UserDTO{
			ID:           following.ID,
			Username:     following.Username,
			FavoriteTeam: favoriteTeam,
		},
		FollowStatsDTO: followStats,
	}

	return &output, nil
}
