package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEditUser(t *testing.T) {
	mockRepo := repository_mocks.NewMockUserRepo()
	mockFootballClient := new(football.MockFootballClient)

	uc := NewEditUserUseCase(mockRepo, mockFootballClient)

	team := football.Team{
		ID:           7633,
		Name:         "Atlético MG",
		Abbreviation: "CAM",
		Color:        "000000",
		Logo:         "https://logo.png",
	}
	mockFootballClient.On("GetTeam", mock.AnythingOfType("int64")).Return(&team, nil)

	existingUser := &domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	user := &domain.User{
		ID:             "uuid2",
		Username:       "modasby2",
		Email:          "modasby2@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	mockRepo.Create(existingUser)
	mockRepo.Create(user)

	t.Run("should edit user successfully", func(t *testing.T) {
		input := EditUserInputDTO{
			UserID:         "uuid2",
			Username:       "modasby3",
			Email:          "modasby3@modasby.com",
			FavoriteTeamID: 7633,
		}

		err := uc.Execute(input)

		assert.NoError(t, err)

		user, err = mockRepo.FindOneByIdOrUsername("modasby3")

		assert.NoError(t, err)

		assert.Equal(t, "modasby3", user.Username)
		assert.Equal(t, "modasby3@modasby.com", user.Email)
		assert.Equal(t, int64(7633), user.FavoriteTeamID)
	})

	t.Run("should return error when email already exists", func(t *testing.T) {
		input := EditUserInputDTO{
			UserID: "uuid2",
			Email:  existingUser.Email,
		}

		err := uc.Execute(input)

		assert.Error(t, err)
	})

	t.Run("should return error when username already exists", func(t *testing.T) {
		input := EditUserInputDTO{
			UserID:   "uuid2",
			Username: existingUser.Username,
		}

		err := uc.Execute(input)

		assert.Error(t, err)
	})

	mockFootballClient.AssertExpectations(t)
}
