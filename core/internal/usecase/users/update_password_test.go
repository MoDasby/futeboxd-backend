package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePassword(t *testing.T) {
	mockUserRepo := repository_mocks.NewMockUserRepo()

	uc := NewUpdatePasswordUsecase(mockUserRepo)

	user := &domain.User{
		ID:             "uuid",
		Name:           "giulliano",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	user.HashPassword()

	mockUserRepo.Create(user)

	input := UpdatePasswordInputDTO{
		User:            user,
		CurrentPassword: "123456",
		NewPassword:     "1234567",
	}

	err := uc.Execute(input)

	assert.NoError(t, err)

	user, _ = mockUserRepo.FindOneByIdOrUsername(user.ID)

	err = user.CheckPassword("1234567")

	assert.NoError(t, err)
}
