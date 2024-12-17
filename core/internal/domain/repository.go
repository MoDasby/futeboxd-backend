package domain

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
)

type UserRepository interface {
	FindOneByCredential(ctx context.Context, credential string) (*User, error)
	Create(ctx context.Context, user *User) error
	Exists(ctx context.Context, username, email string) (bool, error)
	FindOneByIdOrUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	SaveRecoverToken(ctx context.Context, recover *Recover) error
	CheckRecoverToken(ctx context.Context, token string) (*Recover, error)
}

type SessionRepository interface {
	FindOneByToken(ctx context.Context, token string) (*Session, error)
	Create(ctx context.Context, session *Session) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	Delete(ctx context.Context, reviewID int64) error
	ExistsByID(ctx context.Context, reviewID int64) (bool, error)
	ListFeed(ctx context.Context, requesterID string, page *pagination.Page) ([]Review, error)
	ListAll(ctx context.Context, requesterID, where string, params []any, page *pagination.Page) ([]Review, error)
	ListTrendingMatches(ctx context.Context, page *pagination.Page) ([]int64, error)
	Search(ctx context.Context, requesterID, term string, page *pagination.Page) ([]Review, error)
	FindOneByID(ctx context.Context, requesterID string, reviewID int64) (*Review, error)
	Like(ctx context.Context, requesterID string, reviewID int64) error
	Unlike(ctx context.Context, requesterID string, reviewID int64) error
	IsLiked(ctx context.Context, requesterID string, reviewID int64) (bool, error)
}

type CommentsRepository interface {
	Create(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, commentID int64) error
	ExistsByID(ctx context.Context, commentID int64) (bool, error)
	FindOneByID(ctx context.Context, requesterID string, commentID int64) (*Comment, error)
	ListByReview(ctx context.Context, requesterID string, reviewID int64, page *pagination.Page) ([]Comment, error)
	Like(ctx context.Context, requesterID string, commentID int64) error
	Unlike(ctx context.Context, requesterID string, commentID int64) error
	IsLiked(ctx context.Context, requesterID string, commentID int64) (bool, error)
}

type ProfileRepository interface {
	FindOneByUsername(ctx context.Context, username, requesterID string) (*Profile, error)
	Search(ctx context.Context, requesterID, term string, page *pagination.Page) ([]Profile, error)
	Follow(ctx context.Context, followerID, followingID string) error
	Unfollow(ctx context.Context, followerID, followingID string) error
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
}

type MatchRepository interface {
	ListPopularMatches(ctx context.Context, page *pagination.Page) ([]int64, error)
}
