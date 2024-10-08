package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		email          string
		password       string
		favoriteTeamID int64
		expectedError  bool
	}{
		{
			name:           "valid user",
			username:       "modasby",
			email:          "giulliano@email.com",
			password:       "123456",
			favoriteTeamID: 7632,
			expectedError:  false,
		},
		{
			name:           "empty username",
			username:       "",
			email:          "giulliano@email.com",
			password:       "123456",
			favoriteTeamID: 7632,
			expectedError:  true,
		},
		{
			name:           "invalid email",
			username:       "modasby",
			email:          "giullianoemail.com",
			password:       "123456",
			favoriteTeamID: 7632,
			expectedError:  true,
		},
		{
			name:           "invalid password length",
			username:       "modasby",
			email:          "giulliano@email.com",
			password:       "1234",
			favoriteTeamID: 7632,
			expectedError:  true,
		},
		{
			name:           "empty favorite team ID",
			username:       "modasby",
			email:          "giulliano@email.com",
			password:       "123456",
			favoriteTeamID: 0,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.username, tt.email, tt.password, tt.favoriteTeamID)

			if tt.expectedError {
				assert.NotNil(t, err, "expected error, got nil")
				assert.Nil(t, user, "expected nil user, got user")
			} else {
				assert.Nil(t, err, "expected no error, got error")
				assert.NotNil(t, user, "expected user, got nil")

				assert.Equal(t, tt.username, user.Username)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.favoriteTeamID, user.FavoriteTeamID)
			}
		})
	}
}

func TestHash(t *testing.T) {
	password := "123456"
	user, err := NewUser("modasby", "giulliano@email.com", password, 0)

	assert.Nil(t, err)

	err = user.CheckPassword(password)

	assert.Nil(t, err)

	err = user.CheckPassword("1234567")

	assert.NotNil(t, err)
}
