package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteComment(t *testing.T) {
	commentsRepo := repository_mocks.NewMockCommentsRepo()

	comment := &domain.Comment{
		ID: 1,
		Author: &domain.User{
			ID:       "uuid",
			Username: "modasby",
		},
		ParentID: 12,
		Content:  "partida bacana",
	}

	commentsRepo.Create(comment)

	uc := NewDeleteCommentUsecase(commentsRepo)

	t.Run("should return error when creating non-existing comment", func(t *testing.T) {
		err := uc.Execute("uuid", 100)

		assert.Error(t, err)
	})

	t.Run("should return error when user doesn't have authorization to delete comment", func(t *testing.T) {
		err := uc.Execute("otheruser", 1)

		assert.Error(t, err)
	})

	t.Run("should delete comment", func(t *testing.T) {
		err := uc.Execute("uuid", comment.ID)

		assert.NoError(t, err)
	})
}
