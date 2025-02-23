package review

import (
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
)

type Review struct {
	ID            int
	Author        *user.User
	Rate          int
	Description   string
	MatchID       int64
	HomeTeamID    int64
	AwayTeamID    int64
	Likes         int
	CommentsCount int
	IsLiked       bool
	CreatedAt     time.Time
}

func NewReview(
	author *user.User,
	rate int,
	matchID int64,
	description string,
	homeTeamID int64,
	awayTeamID int64,
) (*Review, error) {
	review := &Review{
		Author:      author,
		Rate:        rate,
		Description: description,
		MatchID:     matchID,
		HomeTeamID:  homeTeamID,
		AwayTeamID:  awayTeamID,
	}

	if err := review.Validate(); err != nil {
		return nil, err
	}

	return review, nil
}

func (r *Review) Validate() error {
	if r.Rate < 0 || r.Rate > 5 {
		return &errors.HTTPErr{
			Msg:        "A nota deve estar no intervalo entre 0 e 5",
			Code:       http.StatusBadRequest,
			Context:    "REVIEW:DOMAIN:INVALID_RATE",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  r.Author.ID,
			Timestamp:  time.Now().UTC(),
		}
	}

	return nil
}
