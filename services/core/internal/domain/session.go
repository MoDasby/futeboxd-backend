package domain

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"
)

const (
	DEFAULT_EXPIRATION time.Duration = 24 * 31 * time.Hour
)

type Session struct {
	ID        string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UserID    string
}

func generateToken() (string, error) {
	bytes := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func NewSession(userID string) (*Session, error) {
	token, err := generateToken()
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
