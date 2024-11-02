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

	user := &domain.User{
		Username:       "username",
		Email:          "email@email.com",
		Password:       "password",
		FavoriteTeamID: 7632,
	}

	input := EditUserInputDTO{
		UserID:         user.ID,
		Username:       user.Username,
		Email:          user.Email,
		FavoriteTeamID: user.FavoriteTeamID,
	}

	team := football.Team{
		ID:           7632,
		Name:         "Atlético MG",
		Abbreviation: "CAM",
		Color:        "000000",
		Logo:         "https://logo.png",
	}
	mockFootballClient.On("GetTeam", mock.AnythingOfType("int64")).Return(&team, nil)

	mockRepo.Create(user)

	err := uc.Execute(input)

	assert.Nil(t, err)
}
