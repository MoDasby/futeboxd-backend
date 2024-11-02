package repository_mocks

import (
	"errors"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type MockUserRepo struct {
	users []domain.User
}

func NewMockUserRepo() domain.UserRepository {
	return &MockUserRepo{
		users: make([]domain.User, 0),
	}
}

func (m *MockUserRepo) FindOneByCredential(credential string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Username == credential || user.Email == credential {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepo) Create(user *domain.User) (*domain.User, error) {
	m.users = append(m.users, *user)
	return user, nil
}

func (m *MockUserRepo) Exists(username, email string) (bool, error) {
	for _, user := range m.users {
		if user.Username == username || user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockUserRepo) FindOneByIdOrUsername(username string) (*domain.User, error) {
	for _, user := range m.users {
		if user.ID == username || user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepo) Update(user *domain.User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return errors.New("user not found")
}
