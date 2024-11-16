package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
)

type FindProfileUsecase struct {
	profileRepository domain.ProfileRepository
	footballClient    football.Client
}

func NewFindProfileUsecase(
	profileRepository domain.ProfileRepository,
	footballClient football.Client,
) *FindProfileUsecase {
	return &FindProfileUsecase{
		profileRepository: profileRepository,
		footballClient:    footballClient,
	}
}

func (uc *FindProfileUsecase) Execute(requester *domain.User, username string) (*dto.ProfileDTO, error) {
	profile, err := uc.profileRepository.FindOneByUsername(username, requester.ID)
	if err != nil {
		return nil, err
	}

	favoriteTeam, err := uc.footballClient.GetTeam(profile.FavoriteTeam)
	if err != nil {
		return nil, err
	}

	output := dto.ProfileDTO{
		UserDTO: &dto.UserDTO{
			ID:           profile.UserID,
			Name:         profile.Name,
			Bio:          profile.Bio,
			Username:     profile.Username,
			FavoriteTeam: favoriteTeam,
		},
		FollowersCount: profile.FollowersCount,
		FollowingCount: profile.FollowingCount,
		Following:      profile.IsFollowing,
		Self:           requester.Username == profile.Username,
	}

	return &output, nil
}
