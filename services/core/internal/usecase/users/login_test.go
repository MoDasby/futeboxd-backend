package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	mockUserRepo := repository_mocks.NewMockUserRepo()
	mockSessionRepo := repository_mocks.NewMockSessionRepo()

	uc := NewLoginUseCase(mockSessionRepo, mockUserRepo)

	input := LoginInput{
		Credential: "modasby",
		Password:   "123456",
	}

	user := domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	user.HashPassword()

	mockUserRepo.Create(&user)

	output, err := uc.Execute(&input)

	assert.Nil(t, err)
	assert.NotEmpty(t, output.Token)
	assert.NotEmpty(t, output.ExpiresAt)
	assert.NotEmpty(t, output.CreatedAt)
}
