package repository_mocks

import (
	"errors"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type MockProfileRepo struct {
	profiles []domain.Profile
}

func NewMockProfileRepo() *MockProfileRepo {
	return &MockProfileRepo{
		profiles: make([]domain.Profile, 0),
	}
}

func (m *MockProfileRepo) FindOneByUsername(username, requesterID string) (*domain.Profile, error) {
	for _, profile := range m.profiles {
		if profile.Username == username {
			profile.IsFollowing = false
			return &profile, nil
		}
	}
	return nil, errors.New("profile not found")
}

func (m *MockProfileRepo) AddProfile(profile domain.Profile) {
	m.profiles = append(m.profiles, profile)
}

func (m *MockProfileRepo) Search(requesterID, term string, page *pagination.Page) ([]domain.Profile, error) {
	return nil, nil
}
