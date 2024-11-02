package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindProfile(t *testing.T) {
	mockRepo := new(repository.MockProfileRepository)
	mockFootballClient := new(football.MockFootballClient)

	uc := NewFindProfileUsecase(mockRepo, mockFootballClient)

	username := "modasby2"
	requester := domain.User{
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	mockRepo.On("FindOneByUsername", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&domain.Profile{
		UserID:         "wiufwbfuwbuw",
		Username:       "modasby2",
		FavoriteTeam:   7632,
		FollowersCount: 1937,
		FollowingCount: 1938,
		IsFollowing:    false,
	}, nil)

	team := football.Team{
		ID:           7632,
		Name:         "Atlético MG",
		Abbreviation: "CAM",
		Color:        "000000",
		Logo:         "https://logo.png",
	}
	mockFootballClient.On("GetTeam", mock.AnythingOfType("int64")).Return(&team, nil)

	profile, err := uc.Execute(&requester, username)

	assert.Nil(t, err)
	assert.Equal(t, "modasby2", profile.Username)
	assert.Equal(t, int64(7632), profile.FavoriteTeam.ID)
	assert.Equal(t, "Atlético MG", profile.FavoriteTeam.Name)
	assert.Equal(t, "CAM", profile.FavoriteTeam.Abbreviation)
	assert.Equal(t, "000000", profile.FavoriteTeam.Color)
	assert.Equal(t, "https://logo.png", profile.FavoriteTeam.Logo)
	assert.Equal(t, 1937, profile.FollowersCount)
	assert.Equal(t, 1938, profile.FollowingCount)
	assert.False(t, profile.Following)

	mockRepo.AssertExpectations(t)
	mockFootballClient.AssertExpectations(t)
}
