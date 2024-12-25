package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/modasby/futeboxd-backend/core/pkg/token"
)

const (
	DEFAULT_EXPIRATION = 24 * 31 * time.Hour
)

type Session struct {
	ID        string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    string
}

func NewAnonymousSession() *Session {
	return &Session{
		ID:     "anonymous",
		UserID: uuid.NewString(),
	}
}

func NewSession(userID string) (*Session, error) {
	token, err := token.Generate()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(DEFAULT_EXPIRATION)

	return &Session{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}, nil
}

func (s *Session) IsValid() bool {
	return s.ExpiresAt.UTC().After(time.Now().UTC())
}

func (s *Session) DefaultExpiration() time.Duration {
	return DEFAULT_EXPIRATION
}
