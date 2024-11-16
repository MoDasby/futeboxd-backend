package repository_mocks

import (
	"errors"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type MockFollowersRepo struct {
	followers map[string]map[string]bool
}

func NewMockFollowersRepo() domain.FollowersRepository {
	return &MockFollowersRepo{
		followers: make(map[string]map[string]bool),
	}
}

func (m *MockFollowersRepo) Follow(followerID, followingID string) error {
	if m.followers[followerID] == nil {
		m.followers[followerID] = make(map[string]bool)
	}
	m.followers[followerID][followingID] = true
	return nil
}

func (m *MockFollowersRepo) Unfollow(followerID, followingID string) error {
	if _, ok := m.followers[followerID]; !ok {
		return errors.New("follower not found")
	}
	if _, ok := m.followers[followerID][followingID]; !ok {
		return errors.New("not following this user")
	}
	delete(m.followers[followerID], followingID)
	return nil
}

func (m *MockFollowersRepo) IsFollowing(followerID, followingID string) (bool, error) {
	if _, ok := m.followers[followerID]; !ok {
		return false, nil
	}
	isFollowing := m.followers[followerID][followingID]
	return isFollowing, nil
}
