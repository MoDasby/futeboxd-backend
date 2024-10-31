package repository

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) FindOneByUsername(username, requesterID string) (*domain.Profile, error) {
	args := m.Called(username, requesterID)

	return args.Get(0).(*domain.Profile), args.Error(1)
}
