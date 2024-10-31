package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogout(t *testing.T) {
	mockSessionRepo := new(repository.MockSessionRepo)

	mockSessionRepo.On("Delete", mock.AnythingOfType("string")).Return(nil)

	uc := NewLogoutUsecase(mockSessionRepo)

	session, _ := domain.NewSession("uuid")

	err := uc.Execute(session)

	assert.Nil(t, err)

	mockSessionRepo.AssertExpectations(t)
}
