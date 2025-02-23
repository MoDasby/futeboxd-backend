package review

import (
	"testing"

	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestReview(t *testing.T) {
	author := &user.User{
		ID:       "uuid",
		Username: "modasby",
	}

	t.Run("should create a review", func(t *testing.T) {
		review, err := NewReview(
			author,
			4,
			76453,
			"que partida bacana",
			7632,
			2022,
		)

		assert.NoError(t, err)
		assert.NotNil(t, review)
	})

	t.Run("check validations", func(t *testing.T) {
		review := &Review{
			Rate:   10,
			Author: author,
		}

		err := review.Validate()

		assert.Error(t, err)

		review.Rate = -10

		err = review.Validate()

		assert.Error(t, err)

		review.Rate = 5

		err = review.Validate()

		assert.NoError(t, err)

		review.Description = "partida bacana"

		err = review.Validate()

		assert.NoError(t, err)
	})
}
