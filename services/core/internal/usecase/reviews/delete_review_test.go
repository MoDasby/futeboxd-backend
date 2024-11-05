package usecase

import (
	"testing"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	repository_mocks "github.com/modasby/futeboxd-api/services/core/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteReview(t *testing.T) {
	reviewRepo := repository_mocks.NewMockReviewRepo()

	t.Run("should delete review", func(t *testing.T) {
		user := domain.User{
			ID:             "uuid",
			Username:       "modasby",
			Email:          "modasby@email.com",
			Password:       "123456",
			FavoriteTeamID: 7632,
		}

		review := domain.Review{
			ID:          1,
			Author:      &user,
			Rate:        2,
			Description: "partida lega",
		}

		reviewRepo.Create(&review)

		uc := NewDeleteReviewUseCase(reviewRepo)

		err := uc.Execute(user.ID, int64(review.ID))

		assert.NoError(t, err)
	})

	t.Run("should fail when user does't have authorization", func(t *testing.T) {
		user := domain.User{
			ID:             "uuid",
			Username:       "modasby",
			Email:          "modasby@email.com",
			Password:       "123456",
			FavoriteTeamID: 7632,
		}

		otherUser := domain.User{
			ID: "uuid2",
		}

		review := domain.Review{
			ID:          1,
			Author:      &user,
			Rate:        2,
			Description: "partida lega",
		}

		reviewRepo.Create(&review)

		uc := NewDeleteReviewUseCase(reviewRepo)

		err := uc.Execute(otherUser.ID, int64(review.ID))

		assert.Error(t, err)
	})
}
