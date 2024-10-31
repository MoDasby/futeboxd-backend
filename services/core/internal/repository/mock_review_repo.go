package repository

import (
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockReviewRepo struct {
	mock.Mock
}

func (m *MockReviewRepo) Create(review *domain.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewRepo) Delete(reviewID int64) error {
	args := m.Called(reviewID)
	return args.Error(0)
}

func (m *MockReviewRepo) ListFeed(requesterID, strategy string, pageSize, pageIndex int) ([]domain.Review, error) {
	args := m.Called(requesterID, strategy, pageSize, pageIndex)
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewRepo) ListAll(pageSize, pageIndex int, userID, team, match string) ([]domain.Review, error) {
	args := m.Called(pageSize, pageIndex, userID, team, match)
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewRepo) FindOneByID(reviewID int64) (*domain.Review, error) {
	args := m.Called(reviewID)
	return args.Get(0).(*domain.Review), args.Error(1)
}
