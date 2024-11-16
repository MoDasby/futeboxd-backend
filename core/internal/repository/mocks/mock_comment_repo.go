package repository_mocks

import (
	"errors"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type MockCommentsRepo struct {
	comments []domain.Comment
}

func NewMockCommentsRepo() domain.CommentsRepository {
	return &MockCommentsRepo{
		comments: make([]domain.Comment, 0),
	}
}

func (m *MockCommentsRepo) Create(comment *domain.Comment) error {
	m.comments = append(m.comments, *comment)
	return nil
}

func (m *MockCommentsRepo) Delete(commentID int64) error {
	for i, comment := range m.comments {
		if comment.ID == commentID {
			m.comments = append(m.comments[:i], m.comments[i+1:]...)
			return nil
		}
	}
	return errors.New("comment not found")
}

func (m *MockCommentsRepo) ExistsByID(commentID int64) (bool, error) {
	for _, comment := range m.comments {
		if comment.ID == commentID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockCommentsRepo) FindOneByID(commentID int64) (*domain.Comment, error) {
	for _, comment := range m.comments {
		if comment.ID == commentID {
			return &comment, nil
		}
	}

	return nil, errors.New("commment not found")
}

func (m *MockCommentsRepo) ListByReview(requesterID string, reviewID int64, page *pagination.Page) ([]domain.Comment, error) {
	start := (page.Index - 1) * page.Size
	end := start + page.Size

	var filteredComments []domain.Comment
	for _, comment := range m.comments {
		if comment.ParentID == reviewID {
			filteredComments = append(filteredComments, comment)
		}
	}

	if start > len(filteredComments) {
		return []domain.Comment{}, nil
	}
	if end > len(filteredComments) {
		end = len(filteredComments)
	}
	return filteredComments[start:end], nil
}

func (m *MockCommentsRepo) Like(requesterID string, commentID int64) error {
	for i, comment := range m.comments {
		if comment.ID == commentID {
			comment.LikeCount++
			comment.IsLiked = true
			m.comments[i] = comment
			return nil
		}
	}
	return errors.New("comment not found")
}

func (m *MockCommentsRepo) Unlike(requesterID string, commentID int64) error {
	for i, comment := range m.comments {
		if comment.ID == commentID {
			if comment.LikeCount > 0 {
				comment.LikeCount--
			}
			comment.IsLiked = false
			m.comments[i] = comment
			return nil
		}
	}
	return errors.New("comment not found")
}

func (m *MockCommentsRepo) IsLiked(requesterID string, commentID int64) (bool, error) {
	for _, comment := range m.comments {
		if comment.ID == commentID {
			return comment.IsLiked, nil
		}
	}
	return false, errors.New("comment not found")
}
