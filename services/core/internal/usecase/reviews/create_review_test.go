package usecase

/* import (
	"errors"
	"testing"

	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateReviewUseCase_Execute(t *testing.T) {
	user := domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	t.Run("should create review successfully", func(t *testing.T) {
		mockFootballClient := new(football.MockFootballClient)
		mockReviewRepo := new(repository.MockReviewRepo)

		match := &football.Match{
			ID:             12121,
			HomeCompetitor: &football.Competitor{Team: &football.Team{ID: 7632}},
			AwayCompetitor: &football.Competitor{Team: &football.Team{ID: 7632}},
		}
		mockFootballClient.On("GetMatch", mock.AnythingOfType("int64")).Return(match, nil)
		mockReviewRepo.On("Create", mock.AnythingOfType("*domain.Review")).Return(nil)

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
		mockReviewRepo.AssertExpectations(t)
	})

	t.Run("should return error if GetMatch fails", func(t *testing.T) {
		mockFootballClient := new(football.MockFootballClient)
		mockReviewRepo := new(repository.MockReviewRepo)

		mockFootballClient.On("GetMatch", mock.AnythingOfType("int64")).Return(nil, errors.New("match not found"))

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

	t.Run("should return error if Create in repository fails", func(t *testing.T) {
		mockFootballClient := new(football.MockFootballClient)
		mockReviewRepo := new(repository.MockReviewRepo)

		match := &football.Match{
			ID:             114125,
			HomeCompetitor: &football.Competitor{Team: &football.Team{ID: 7632}},
			AwayCompetitor: &football.Competitor{Team: &football.Team{ID: 2022}},
		}
		mockFootballClient.On("GetMatch", mock.AnythingOfType("int64")).Return(match, nil)
		mockReviewRepo.On("Create", mock.AnythingOfType("*domain.Review")).Return(errors.New("db error"))

		uc := NewCreateReviewUseCase(mockReviewRepo, mockFootballClient)

		input := CreateReviewInputDTO{
			Author:      &user,
			Rate:        5,
			MatchID:     3535,
			Description: "Great match!",
		}

		err := uc.Execute(input)

		assert.EqualError(t, err, "db error")
		mockFootballClient.AssertExpectations(t)
		mockReviewRepo.AssertExpectations(t)
	})
} */
