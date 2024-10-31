package usecase

import (
	"testing"
	"time"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	mockUserRepo := new(repository.MockUserRepo)
	mockSessionRepo := new(repository.MockSessionRepo)

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

	session := domain.Session{
		Token:     "generatedToken",
		ExpiresAt: time.Now().Add(domain.DEFAULT_EXPIRATION),
		CreatedAt: time.Now(),
		UserID:    "uuid",
	}

	mockUserRepo.On("FindOneByCredential", mock.AnythingOfType("string")).Return(&user, nil)

	mockSessionRepo.On("Create", mock.AnythingOfType("*domain.Session")).Return(&session, nil)

	output, err := uc.Execute(&input)

	assert.Nil(t, err)
	assert.Equal(t, session.Token, output.Token)
	assert.Equal(t, session.ExpiresAt, output.ExpiresAt)
	assert.Equal(t, session.CreatedAt, output.CreatedAt)

	mockSessionRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
