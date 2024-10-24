package domain

import "github.com/modasby/futeboxd-api/pkg/pagination"

type SessionRepository interface {
	FindOneByToken(token string) (*Session, error)
	Create(session *Session) (*Session, error)
	Delete(id string) error
}

type ReviewRepository interface {
	Create(review *Review) error
	Delete(reviewID int64) error
	ListFeed(requesterID, strategy string, pageSize, pageIndex int) ([]Review, error)
	ListAll(pageSize, pageIndex int, userID, team, match string) ([]Review, error)
	FindOneByID(reviewID int64) (*Review, error)
}

type UserRepository interface {
	FindOneByCredential(credential string) (*User, error)
	Create(user *User) (*User, error)
	Exists(username, email string) (bool, error)
	FindOneByIdOrUsername(username string) (*User, error)
	Update(user *User) error
}

type FollowersRepository interface {
	Follow(followerID, followingID string) error
	Unfollow(followerID, followingID string) error
	IsFollowing(followerID, followingID string) (bool, error)
}

type CommentsRepository interface {
	Create(comment *Comment) error
	Delete(commentID int64) error
	ListByReview(reviewID int64, page *pagination.Page) ([]Comment, error)
}
