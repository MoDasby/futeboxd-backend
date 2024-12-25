package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {

	t.Run("should create session", func(t *testing.T) {
		session, err := NewSession("uuid")

		assert.NoError(t, err)

		assert.Len(t, session.Token, 96)
		assert.False(t, session.ExpiresAt.IsZero())
		assert.Equal(t, "uuid", session.UserID)
	})

	t.Run("check validation", func(t *testing.T) {
		session := &Session{
			ExpiresAt: time.Now().Add(DEFAULT_EXPIRATION),
		}

		valid := session.IsValid()

		assert.True(t, valid)

		session.ExpiresAt = time.Now().Add(-1 * time.Hour)

		valid = session.IsValid()

		assert.False(t, valid)
	})
}
