package repository

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindOneByCredential(credential string) (*domain.User, error) {
	args := m.Called(credential)

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Create(user *domain.User) (*domain.User, error) {
	args := m.Called(user)

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Exists(username, email string) (bool, error) {
	args := m.Called(username, email)

	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) FindOneByIdOrUsername(username string) (*domain.User, error) {
	args := m.Called(username)

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Update(user *domain.User) error {
	args := m.Called(user)

	return args.Error(0)
}
