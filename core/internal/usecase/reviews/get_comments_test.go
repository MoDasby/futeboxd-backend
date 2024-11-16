package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetComments(t *testing.T) {
	t.Run("should return empty slice", func(t *testing.T) {
		commentsRepo := repository_mocks.NewMockCommentsRepo()
		reviewsRepo := repository_mocks.NewMockReviewRepo()

		reviewsRepo.Create(&domain.Review{
			ID: 1,
		})

		commentsRepo.Create(&domain.Comment{
			ID: 1,
			Author: &domain.User{
				ID:       "uuid",
				Username: "modasby",
			},
			ParentID: 1,
		})
		uc := NewGetCommentsUsecase(commentsRepo)

		comments, err := uc.Execute("user1", 2, &pagination.Page{Index: 1, Size: 10})

		assert.NoError(t, err)
		assert.Len(t, comments, 0)

		comments, err = uc.Execute("uuid", 1, &pagination.Page{Index: 1, Size: 10})

		assert.NoError(t, err)
		assert.Len(t, comments, 1)
	})
}
