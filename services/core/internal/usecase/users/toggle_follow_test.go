package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestToggleFollow(t *testing.T) {
	mockUserRepo := repository_mocks.NewMockUserRepo()
	mockFollowersRepo := repository_mocks.NewMockFollowersRepo()

	uc := NewToggleFollowUsecase(mockUserRepo, mockFollowersRepo)

	follower := &domain.User{
		ID:             "uuid2",
		Username:       "modasby2",
		Email:          "modasby2@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	following := &domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	mockUserRepo.Create(follower)
	mockUserRepo.Create(following)

	output, err := uc.Execute(follower.ID, following.ID)

	assert.NoError(t, err)
	assert.True(t, output.Follow)

	output, err = uc.Execute(follower.ID, following.ID)

	assert.NoError(t, err)
	assert.True(t, output.Follow)
}
