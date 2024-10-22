package usecase

import "github.com/modasby/futeboxd-api/services/core/internal/domain"

type UnfollowUserUseCase struct {
	followersRepo domain.FollowersRepository
	userRepo      domain.UserRepository
}

func NewUnfollowUserUseCase(
	followerRepo domain.FollowersRepository,
	userRepo domain.UserRepository,
) *UnfollowUserUseCase {
	return &UnfollowUserUseCase{
		followersRepo: followerRepo,
		userRepo:      userRepo,
	}
}

func (uc *UnfollowUserUseCase) Execute(followerID, followingUsername string) error {
	following, err := uc.userRepo.FindOneByIdOrUsername(followingUsername)
	if err != nil {
		return err
	}

	isFollowing, err := uc.followersRepo.IsFollowing(followerID, following.ID)
	if err != nil {
		return err
	}

	if followerID == following.ID || !isFollowing {
		return nil
	}

	if err := uc.followersRepo.Unfollow(followerID, following.ID); err != nil {
		return err
	}
	return nil
}
