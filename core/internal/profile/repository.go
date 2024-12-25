package profile

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type Repository interface {
	FindOneByUsername(ctx context.Context, username, requesterID string) (*Profile, error)
	Search(ctx context.Context, requesterID, term string, page *pagination.Page) ([]Profile, error)
	Follow(ctx context.Context, followerID, followingID string) error
	Unfollow(ctx context.Context, followerID, followingID string) error
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
}
