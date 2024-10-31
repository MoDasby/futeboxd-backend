package domain

import "github.com/modasby/futeboxd-api/pkg/pagination"

type SessionRepository interface {
	FindOneByToken(token string) (*Session, error)
	Create(session *Session) (*Session, error)
	Update(session *Session) error
	Delete(id string) error
}

type ReviewRepository interface {
	Create(review *Review) error
	Delete(reviewID int64) error
	ExistsByID(reviewID int64) (bool, error)
	ListFeed(requesterID, strategy string, page *pagination.Page) ([]Review, error)
	ListAll(requesterID, userID, team, match string, page *pagination.Page) ([]Review, error)
	FindOneByID(requesterID string, reviewID int64) (*Review, error)
	Like(requesterID string, reviewID int64) error
	Unlike(requesterID string, reviewID int64) error
	IsLiked(requesterID string, reviewID int64) (bool, error)
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
	ExistsByID(commentID int64) (bool, error)
	ListByReview(requesterID string, reviewID int64, page *pagination.Page) ([]Comment, error)
	Like(requesterID string, commentID int64) error
	Unlike(requesterID string, commentID int64) error
	IsLiked(requesterID string, commentID int64) (bool, error)
}

type ProfileRepository interface {
	FindOneByUsername(username, requesterID string) (*Profile, error)
}
