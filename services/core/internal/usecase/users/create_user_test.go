package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/pkg/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserUsecase(t *testing.T) {
	mockRepo := new(repository.MockUserRepo)
	mockFootballClient := new(football.MockFootballClient)

	createUserUsecase := NewCreateUserUseCase(mockRepo, mockFootballClient)

	user := domain.User{
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	input := dto.UserInputDTO{
		Username:       user.Username,
		Email:          user.Email,
		Password:       user.Password,
		FavoriteTeamID: user.FavoriteTeamID,
	}

	mockRepo.On("Exists", user.Username, user.Email).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(&user, nil)

	team := football.Team{
		ID:           7632,
		Name:         "Atlético MG",
		Abbreviation: "CAM",
		Color:        "000000",
		Logo:         "https://logo.png",
	}
	mockFootballClient.On("GetTeam", user.FavoriteTeamID).Return(&team, nil)

	output, err := createUserUsecase.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, user.Username, output.Username)
	assert.Equal(t, user.Email, output.Email)
	assert.Equal(t, team, *output.FavoriteTeam)

	mockRepo.AssertExpectations(t)
	mockFootballClient.AssertExpectations(t)
}
