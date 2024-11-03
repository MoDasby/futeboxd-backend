package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserUsecase(t *testing.T) {
	mockRepo := repository_mocks.NewMockUserRepo()
	mockFootballClient := new(football.MockFootballClient)

	createUserUsecase := NewCreateUserUseCase(mockRepo, mockFootballClient)

	input := dto.UserInputDTO{
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	t.Run("should create a user sucessfully", func(t *testing.T) {
		team := football.Team{
			ID:           7632,
			Name:         "Atlético MG",
			Abbreviation: "CAM",
			Color:        "000000",
			Logo:         "https://logo.png",
		}
		mockFootballClient.On("GetTeam", input.FavoriteTeamID).Return(&team, nil)

		output, err := createUserUsecase.Execute(input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, input.Username, output.Username)
		assert.Equal(t, input.Email, output.Email)
		assert.Equal(t, team, *output.FavoriteTeam)

		mockFootballClient.AssertExpectations(t)
	})

	t.Run("should return an error when creating a user that already exists", func(t *testing.T) {
		output, err := createUserUsecase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
	})
}
