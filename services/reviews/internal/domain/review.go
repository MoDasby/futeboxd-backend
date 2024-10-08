package domain

import "github.com/modasby/futeboxd-api/pkg/errors"

type Review struct {
	ID          int
	UserID      string
	Rate        int
	Description string
	MatchID     int64
	HomeTeamID  int64
	AwayTeamID  int64
}

func NewReview(userID string, rate int, matchID int64, description string, homeTeamID int64, awayTeamID int64) (*Review, error) {
	review := &Review{
		UserID:      userID,
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
	if r.UserID == "" {
		return errors.NewErrBadRequest("usuário inválido ao criar review")
	}

	if r.Rate <= 0 {
		return errors.NewErrBadRequest("avaliação deve ser pelo menos 1")
	}

	if r.Description == "" {
		return errors.NewErrBadRequest("insira uma descrição")
	}

	return nil
}
