package repository

import (
	"github.com/modasby/futeboxd-api/pkg/pagination"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) Create(comment *domain.Comment) error {
	args := m.Called(comment)

	return args.Error(0)
}

func (m *MockCommentRepo) Delete(commentID int64) error {
	args := m.Called(commentID)

	return args.Error(0)
}

func (m *MockCommentRepo) ListByReview(reviewID int64, page *pagination.Page) ([]domain.Comment, error) {
	args := m.Called(reviewID, page)

	return args.Get(0).([]domain.Comment), args.Error(1)
}
