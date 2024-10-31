package repository

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockSessionRepo struct {
	mock.Mock
}

func (m *MockSessionRepo) FindOneByToken(token string) (*domain.Session, error) {
	args := m.Called(token)

	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepo) Update(session *domain.Session) error {
	args := m.Called(session)

	return args.Error(0)
}

func (m *MockSessionRepo) Create(session *domain.Session) (*domain.Session, error) {
	args := m.Called(session)

	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepo) Delete(id string) error {
	args := m.Called(id)

	return args.Error(0)
}
