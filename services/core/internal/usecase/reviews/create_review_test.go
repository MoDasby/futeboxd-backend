package usecase

import (
	"errors"
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateReviewUseCase_Execute(t *testing.T) {
	user := domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	mockFootballClient := new(football.MockFootballClient)
	mockReviewRepo := repository_mocks.NewMockReviewRepo()

	t.Run("should create review successfully", func(t *testing.T) {

		match := &football.Match{
			ID:             12121,
			HomeCompetitor: football.Competitor{Team: football.Team{ID: 7632}},
			AwayCompetitor: football.Competitor{Team: football.Team{ID: 7632}},
		}
		mockFootballClient.On("GetMatch", int64(24242)).Return(match, nil)

		uc := NewCreateReviewUseCase(mockReviewRepo, mockFootballClient)

		input := CreateReviewInputDTO{
			Author:      &user,
			Rate:        5,
			MatchID:     24242,
			Description: "Great match!",
		}

		err := uc.Execute(input)

		assert.NoError(t, err)
		mockFootballClient.AssertExpectations(t)
	})

	t.Run("should return error if GetMatch fails", func(t *testing.T) {

		mockFootballClient.On("GetMatch", int64(32535)).Return(nil, errors.New("match not found"))

		uc := NewCreateReviewUseCase(mockReviewRepo, mockFootballClient)

		input := CreateReviewInputDTO{
			Author:      &user,
			Rate:        5,
			MatchID:     32535,
			Description: "Great match!",
		}

		err := uc.Execute(input)

		assert.EqualError(t, err, "match not found")
		mockFootballClient.AssertExpectations(t)
	})
}
