package repository

import "github.com/stretchr/testify/mock"

type MockFollowersRepo struct {
	mock.Mock
}

func (m *MockFollowersRepo) Follow(followerID, followingID string) error {
	args := m.Called(followerID, followingID)

	return args.Error(0)
}

func (m *MockFollowersRepo) Unfollow(followerID, followingID string) error {
	args := m.Called(followerID, followingID)

	return args.Error(0)
}

func (m *MockFollowersRepo) IsFollowing(followerID, followingID string) (bool, error) {
	args := m.Called(followerID, followingID)

	return args.Bool(0), args.Error(1)
}
