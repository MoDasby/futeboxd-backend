package auth

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
