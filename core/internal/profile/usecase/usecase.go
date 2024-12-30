package usecase

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/profile"
	"github.com/modasby/futeboxd-backend/core/internal/profile/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type profileUsecases struct {
	profileRepo    profile.Repository
	footballClient football.Client
}

func NewProfileUsecases(
	profileRepo profile.Repository,
	footballClient football.Client,
) profile.ProfileUsecases {
	return &profileUsecases{
		profileRepo:    profileRepo,
		footballClient: footballClient,
	}
}

func (uc *profileUsecases) FindByUsername(ctx context.Context, username string) (*dto.Profile, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := uc.profileRepo.FindOneByUsername(ctx, username, session.UserID)
	if err != nil {
		return nil, err
	}

	favoriteTeam, err := uc.footballClient.GetTeam(profile.FavoriteTeam)
	if err != nil {
		return nil, err
	}

	output := dto.Profile{
		ID:             profile.UserID,
		Name:           null.NewString(profile.Name),
		Bio:            null.NewString(profile.Bio),
		Username:       profile.Username,
		FavoriteTeam:   favoriteTeam,
		FollowersCount: profile.FollowersCount,
		FollowingCount: profile.FollowingCount,
		Following:      profile.IsFollowing,
		Self:           session.UserID == profile.UserID,
	}

	return &output, nil
}

func (uc *profileUsecases) SearchByUsername(ctx context.Context, username string, page *pagination.Page) ([]dto.Profile, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if username == "" {
		return make([]dto.Profile, 0), nil
	}

	profiles, err := uc.profileRepo.Search(ctx, session.UserID, username, page)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Profile, len(profiles))

	for index, profile := range profiles {
		favoriteTeam, err := uc.footballClient.GetTeam(profile.FavoriteTeam)
		if err != nil {
			return nil, err
		}

		output[index] = dto.Profile{
			ID:             profile.UserID,
			Name:           null.NewString(profile.Name),
			Bio:            null.NewString(profile.Bio),
			Username:       profile.Username,
			FavoriteTeam:   favoriteTeam,
			FollowersCount: profile.FollowersCount,
			FollowingCount: profile.FollowingCount,
			Following:      profile.IsFollowing,
			Self:           session.UserID == profile.UserID,
		}
	}

	return output, nil
}

func (uc *profileUsecases) ToggleFollow(ctx context.Context, usernameToFollow string) (*dto.FollowStats, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	userToFollow, err := uc.profileRepo.FindOneByUsername(ctx, usernameToFollow, session.UserID)
	if err != nil {
		return nil, err
	}

	if session.UserID == userToFollow.UserID {
		return &dto.FollowStats{Following: false}, nil
	}

	if userToFollow.IsFollowing {
		err := uc.profileRepo.Unfollow(ctx, session.UserID, userToFollow.UserID)
		if err != nil {
			return nil, err
		}

		return &dto.FollowStats{Following: false}, nil
	}

	if err := uc.profileRepo.Follow(ctx, session.UserID, userToFollow.UserID); err != nil {
		return nil, err
	}

	return &dto.FollowStats{Following: true}, nil
}
