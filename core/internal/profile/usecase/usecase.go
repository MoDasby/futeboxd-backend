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

	return uc.toDto(profile, session.UserID)
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

	return uc.toDtoList(profiles, session.UserID)
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

func (uc *profileUsecases) ListFollowers(ctx context.Context, username string, page *pagination.Page) ([]dto.Profile, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profiles, err := uc.profileRepo.ListFollowers(ctx, session.UserID, username, page)
	if err != nil {
		return nil, err
	}

	return uc.toDtoList(profiles, session.UserID)
}

func (uc *profileUsecases) ListFollowing(ctx context.Context, username string, page *pagination.Page) ([]dto.Profile, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profiles, err := uc.profileRepo.ListFollowing(ctx, session.UserID, username, page)
	if err != nil {
		return nil, err
	}

	return uc.toDtoList(profiles, session.UserID)
}

func (uc *profileUsecases) ListPopularProfiles(ctx context.Context, page *pagination.Page) ([]dto.Profile, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profiles, err := uc.profileRepo.ListPopularProfiles(ctx, session.UserID, page)
	if err != nil {
		return nil, err
	}

	return uc.toDtoList(profiles, session.UserID)
}

func (uc *profileUsecases) toDto(profile *profile.Profile, requesterID string) (*dto.Profile, error) {
	var team *football.Team
	var err error

	if profile.FavoriteTeam > 0 {
		team, err = uc.footballClient.GetTeam(profile.FavoriteTeam)
		if err != nil {
			return nil, err
		}
	}

	profileDto := dto.Profile{
		ID:             profile.UserID,
		Name:           null.String(profile.Name),
		Bio:            null.String(profile.Bio),
		ProfilePicture: profile.ProfilePicture,
		Username:       profile.Username,
		FavoriteTeam:   team,
		Stats: &dto.ProfileStats{
			FollowersCount: profile.FollowersCount,
			FollowingCount: profile.FollowingCount,
			Following:      profile.IsFollowing,
			ReviewsCount:   profile.ReviewsCount,
			Self:           profile.UserID == requesterID,
		},
	}

	return &profileDto, nil
}

func (uc *profileUsecases) toDtoList(profiles []profile.Profile, requesterID string) ([]dto.Profile, error) {
	output := make([]dto.Profile, len(profiles))

	for i := range output {
		profile := profiles[i]

		profileDto, err := uc.toDto(&profile, requesterID)
		if err != nil {
			return nil, err
		}

		output[i] = *profileDto
	}

	return output, nil
}
