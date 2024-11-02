package domain

import (
	"time"

	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type Review struct {
	ID          int
	Author      *User
	Rate        int
	Description string
	MatchID     int64
	HomeTeamID  int64
	AwayTeamID  int64
	Likes       int
	IsLiked     bool
	CreatedAt   time.Time
}

func NewReview(
	author *User,
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
		return errors.NewHTTPErr(
			"a nota deve estar no intervalo entre 0 e 5",
			400,
			"DOMAIN:REVIEW:INVALID_RATE",
		)
	}

	if r.Description == "" {
		return errors.NewHTTPErr(
			"insira uma descrição",
			400,
			"DOMAIN:REVIEW:VALIDATE:INVALID_DESCRIPTION",
		)
	}

	return nil
}
