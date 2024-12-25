package comment

import (
	"testing"

	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestComment(t *testing.T) {
	author := &user.User{
		ID:       "uuid",
		Username: "modasby",
	}

	t.Run("should create a comment", func(t *testing.T) {
		comment, err := NewComment(author, 1, "partida bacana")

		assert.NoError(t, err)
		assert.Equal(t, comment.Author, author)
		assert.Equal(t, comment.Content, "partida bacana")
	})

	t.Run("should return an error when comment is invalid", func(t *testing.T) {
		comment, err := NewComment(author, 1, "")

		assert.Error(t, err)
		assert.Nil(t, comment)

		comment, err = NewComment(author, 0, "partida legal")

		assert.Error(t, err)
		assert.Nil(t, comment)
	})
}
