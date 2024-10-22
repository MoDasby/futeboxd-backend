package domain

import (
	"time"
)

type Session struct {
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    string
}

func NewSession(token string, userID string, expiresAt time.Time) *Session {
	return &Session{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}
}
