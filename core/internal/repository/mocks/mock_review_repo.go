package repository_mocks

import (
	"errors"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type MockReviewRepo struct {
	reviews []domain.Review
}

func NewMockReviewRepo() *MockReviewRepo {
	return &MockReviewRepo{
		reviews: make([]domain.Review, 0),
	}
}

func (m *MockReviewRepo) Search(requesterID, term string, page *pagination.Page) ([]domain.Review, error) {
	return nil, nil
}

func (m *MockReviewRepo) ListTrendingMatches(page *pagination.Page) ([]int64, error) {
	return nil, nil // TODO implementar isso
}

func (m *MockReviewRepo) Create(review *domain.Review) error {
	m.reviews = append(m.reviews, *review)
	return nil
}

func (m *MockReviewRepo) Delete(reviewID int64) error {
	for i, review := range m.reviews {
		if int64(review.ID) == reviewID {
			m.reviews = append(m.reviews[:i], m.reviews[i+1:]...)
			return nil
		}
	}
	return errors.New("review not found")
}

func (m *MockReviewRepo) ExistsByID(reviewID int64) (bool, error) {
	for _, review := range m.reviews {
		if int64(review.ID) == reviewID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockReviewRepo) ListFeed(requesterID, strategy string, page *pagination.Page) ([]domain.Review, error) {
	start := (page.Index - 1) * page.Size
	end := start + page.Size

	if start > len(m.reviews) {
		return []domain.Review{}, nil
	}
	if end > len(m.reviews) {
		end = len(m.reviews)
	}
	return m.reviews[start:end], nil
}

func (m *MockReviewRepo) ListAll(requesterID, where string, params []any, page *pagination.Page) ([]domain.Review, error) {
	start := page.Index * page.Size
	end := start + page.Size

	if start > len(m.reviews) {
		return []domain.Review{}, nil
	}
	if end > len(m.reviews) {
		end = len(m.reviews)
	}
	return m.reviews[start:end], nil
}

func (m *MockReviewRepo) FindOneByID(requesterID string, reviewID int64) (*domain.Review, error) {
	for _, review := range m.reviews {
		if int64(review.ID) == reviewID {
			return &review, nil
		}
	}
	return nil, errors.New("review not found")
}

func (m *MockReviewRepo) Like(requesterID string, reviewID int64) error {
	for i, review := range m.reviews {
		if int64(review.ID) == reviewID {
			review.Likes++
			review.IsLiked = true
			m.reviews[i] = review
			return nil
		}
	}
	return errors.New("review not found")
}

func (m *MockReviewRepo) Unlike(requesterID string, reviewID int64) error {
	for i, review := range m.reviews {
		if int64(review.ID) == reviewID {
			if review.Likes > 0 {
				review.Likes--
			}
			review.IsLiked = false
			m.reviews[i] = review
			return nil
		}
	}
	return errors.New("review not found")
}

func (m *MockReviewRepo) IsLiked(requesterID string, reviewID int64) (bool, error) {
	for _, review := range m.reviews {
		if int64(review.ID) == reviewID {
			return review.IsLiked, nil
		}
	}
	return false, errors.New("review not found")
}
