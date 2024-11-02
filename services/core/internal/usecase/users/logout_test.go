package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	mockSessionRepo := repository_mocks.NewMockSessionRepo()

	session, _ := domain.NewSession("uuid")

	mockSessionRepo.Create(session)

	uc := NewLogoutUsecase(mockSessionRepo)

	err := uc.Execute(session)

	assert.Nil(t, err)
}
