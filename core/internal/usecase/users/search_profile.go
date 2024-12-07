package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/json/null"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type SearchProfile struct {
	profileRepo    domain.ProfileRepository
	footballClient football.Client
}

func NewSearchProfile(
	profileRepo domain.ProfileRepository,
	footballClient football.Client,
) *SearchProfile {
	return &SearchProfile{profileRepo: profileRepo, footballClient: footballClient}
}

func (uc *SearchProfile) Execute(requester *domain.User, term string, page *pagination.Page) ([]dto.ProfileDTO, error) {
	if term == "" {
		return make([]dto.ProfileDTO, 0), nil
	}

	profiles, err := uc.profileRepo.Search(requester.ID, term, page)
	if err != nil {
		return nil, err
	}

	output := make([]dto.ProfileDTO, len(profiles))

	for index, profile := range profiles {
		favoriteTeam, err := uc.footballClient.GetTeam(profile.FavoriteTeam)
		if err != nil {
			return nil, err
		}

		output[index] = dto.ProfileDTO{
			ID:             profile.UserID,
			Name:           null.NewString(profile.Name),
			Bio:            null.NewString(profile.Bio),
			Username:       profile.Username,
			FavoriteTeam:   favoriteTeam,
			FollowersCount: profile.FollowersCount,
			FollowingCount: profile.FollowingCount,
			Following:      profile.IsFollowing,
			Self:           requester.Username == profile.Username,
		}
	}

	return output, nil
}
