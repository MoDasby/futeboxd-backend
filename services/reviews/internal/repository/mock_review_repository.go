package repository

import (
	"github.com/modasby/futeboxd-api/services/reviews/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) AddReview(review *domain.Review) error {
	args := m.Called(review)

	return args.Error(0)
}

func (m *MockReviewRepository) ListReviews() ([]*domain.Review, error) {
	args := m.Called()

	return args.Get(1).([]*domain.Review), args.Error(1)
}
