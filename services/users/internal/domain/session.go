package domain

import (
	"time"
)

type Session struct {
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Owner     *User
}

func NewSession(token string, ownerID *User, expiresAt time.Time) *Session {
	return &Session{
		Token:     token,
		ExpiresAt: expiresAt,
		Owner:     ownerID,
	}
}
