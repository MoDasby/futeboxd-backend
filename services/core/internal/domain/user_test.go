package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("should create a valid user", func(t *testing.T) {
		user, err := NewUser("validUser", "user@example.com", "securepassword", 1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "validUser", user.Username)
		assert.Equal(t, "user@example.com", user.Email)
		assert.NotEmpty(t, user.Password)
	})

	t.Run("should return an error for invalid username", func(t *testing.T) {
		user, err := NewUser("invalid user", "user@example.com", "securepassword", 1)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return an error for invalid email", func(t *testing.T) {
		user, err := NewUser("validUser", "invalidemail", "securepassword", 1)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return an error for short password", func(t *testing.T) {
		user, err := NewUser("validUser", "user@example.com", "123", 1)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Run("should update password successfully", func(t *testing.T) {
		user, _ := NewUser("validUser", "user@example.com", "securepassword", 1)

		err := user.UpdatePassword("securepassword", "newsecurepassword")

		assert.NoError(t, err)
		assert.NoError(t, user.CheckPassword("newsecurepassword"))
	})

	t.Run("should return an error for incorrect current password", func(t *testing.T) {
		user, _ := NewUser("validUser", "user@example.com", "securepassword", 1)

		err := user.UpdatePassword("wrongpassword", "newsecurepassword")

		assert.Error(t, err)
	})

	t.Run("should return an error for same new password", func(t *testing.T) {
		user, _ := NewUser("validUser", "user@example.com", "securepassword", 1)

		err := user.UpdatePassword("securepassword", "securepassword")

		assert.Error(t, err)
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("should return no error for correct password", func(t *testing.T) {
		user, _ := NewUser("validUser", "user@example.com", "securepassword", 1)

		err := user.CheckPassword("securepassword")

		assert.NoError(t, err)
	})

	t.Run("should return an error for incorrect password", func(t *testing.T) {
		user, _ := NewUser("validUser", "user@example.com", "securepassword", 1)

		err := user.CheckPassword("wrongpassword")

		assert.Error(t, err)
	})
}

func TestValidate(t *testing.T) {
	t.Run("should return no error for valid user", func(t *testing.T) {
		user := &User{
			Username:       "validUser",
			Email:          "user@example.com",
			Password:       "securepassword",
			FavoriteTeamID: 1,
		}

		err := user.Validate()

		assert.NoError(t, err)
	})

	t.Run("should return an error for username with space", func(t *testing.T) {
		user := &User{
			Username:       "invalid user",
			Email:          "user@example.com",
			Password:       "securepassword",
			FavoriteTeamID: 1,
		}

		err := user.Validate()

		assert.Error(t, err)
	})

	t.Run("should return an error for invalid email", func(t *testing.T) {
		user := &User{
			Username:       "validUser",
			Email:          "invalidemail",
			Password:       "securepassword",
			FavoriteTeamID: 1,
		}

		err := user.Validate()

		assert.Error(t, err)
	})

	t.Run("should return an error for short password", func(t *testing.T) {
		user := &User{
			Username:       "validUser",
			Email:          "user@example.com",
			Password:       "123",
			FavoriteTeamID: 1,
		}

		err := user.Validate()

		assert.Error(t, err)
	})
}
