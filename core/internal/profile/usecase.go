package profile

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/profile/dto"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type ProfileUsecases interface {
	FindByUsername(ctx context.Context, username string) (*dto.Profile, error)
	SearchByUsername(ctx context.Context, username string, page *pagination.Page) ([]dto.Profile, error)
	ToggleFollow(ctx context.Context, usernameToFollow string) (*dto.FollowStats, error)
}
