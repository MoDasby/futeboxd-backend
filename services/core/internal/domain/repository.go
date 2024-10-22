package domain

type SessionRepository interface {
	FindOneByToken(token string) (*Session, error)
	AddSession(session *Session) (*Session, error)
	DeleteSession(id string) error
}

type ReviewRepository interface {
	AddReview(review *Review) error
	DeleteReview(reviewID int64) error
	ListFeed(requesterID, strategy string, pageSize, pageIndex int) ([]Review, error)
	ListReviews(pageSize, pageIndex int, userID, team, match string) ([]Review, error)
	FindReviewByID(reviewID int64) (*Review, error)
}

type UserRepository interface {
	FindOneByCredential(credential string) (*User, error)
	AddUser(user *User) (*User, error)
	Exists(username, email string) (bool, error)
	FindOneByIdOrUsername(username string) (*User, error)
	FindBatchByID(ids []string) ([]User, error)
	EditUser(user *User) error
}

type FollowersRepository interface {
	Follow(followerID, followingID string) error
	Unfollow(followerID, followingID string) error
	IsFollowing(followerID, followingID string) (bool, error)
}

type CommentsRepository interface {
	Create(comment Comment) (Comment, error)
	Delete(commentID int64) error
	GetByReview(reviewID int64) ([]Comment, error)
}
