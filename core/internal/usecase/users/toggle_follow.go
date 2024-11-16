package usecase

import "github.com/modasby/futeboxd-api/services/core/internal/domain"

type ToggleFollowUsecase struct {
	userRepo      domain.UserRepository
	followersRepo domain.FollowersRepository
}

func NewToggleFollowUsecase(
	userRepo domain.UserRepository,
	followersRepo domain.FollowersRepository,
) *ToggleFollowUsecase {
	return &ToggleFollowUsecase{
		userRepo:      userRepo,
		followersRepo: followersRepo,
	}
}

type ToggleFollowOutputDTO struct {
	Follow bool `json:"follow"`
}

func (uc *ToggleFollowUsecase) Execute(followerID, followingUsername string) (*ToggleFollowOutputDTO, error) {
	following, err := uc.userRepo.FindOneByIdOrUsername(followingUsername)
	if err != nil {
		return nil, err
	}

	isFollowing, err := uc.followersRepo.IsFollowing(followerID, following.ID)
	if err != nil {
		return nil, err
	}

	if isFollowing {
		err := uc.followersRepo.Unfollow(followerID, following.ID)
		if err != nil {
			return nil, err
		}

		return &ToggleFollowOutputDTO{Follow: false}, nil
	}

	if err := uc.followersRepo.Follow(followerID, following.ID); err != nil {
		return nil, err
	}

	return &ToggleFollowOutputDTO{Follow: true}, nil
}
