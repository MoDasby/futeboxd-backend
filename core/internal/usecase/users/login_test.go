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

	user := domain.User{
		ID:             "uuid",
		Username:       "modasby",
		Email:          "modasby@email.com",
		Password:       "123456",
		FavoriteTeamID: 7632,
	}

	user.HashPassword()

	mockUserRepo.Create(&user)

	t.Run("should login successfully", func(t *testing.T) {
		input := LoginInput{
			Credential: "modasby",
			Password:   "123456",
		}

		output, err := uc.Execute(&input)

		assert.Nil(t, err)
		assert.NotEmpty(t, output.Token)
		assert.False(t, output.ExpiresAt.IsZero())
		assert.False(t, output.CreatedAt.IsZero())
	})

	t.Run("should return error when invalid credential", func(t *testing.T) {
		input := LoginInput{
			Credential: "invalid_credential",
			Password:   "123456",
		}

		output, err := uc.Execute(&input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("should return error when invalid password", func(t *testing.T) {
		input := LoginInput{
			Credential: "modasby",
			Password:   "invalid_pass",
		}

		output, err := uc.Execute(&input)

		assert.Error(t, err)
		assert.Nil(t, output)
	})
}
