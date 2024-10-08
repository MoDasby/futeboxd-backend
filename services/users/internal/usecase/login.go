package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"

	"github.com/modasby/futeboxd-api/services/users/internal/domain"
)

type LoginUseCase struct {
	sessionRepository domain.SessionRepository
	userRepository    domain.UserRepository
}

func NewLoginUseCase(
	sessionRepository domain.SessionRepository,
	userRepository domain.UserRepository,
) *LoginUseCase {
	return &LoginUseCase{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

type LoginInput struct {
	Credential string `json:"credential"`
	Password   string `json:"password"`
}

type LoginOutput struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (uc *LoginUseCase) Execute(input *LoginInput) (*LoginOutput, error) {
	user, err := uc.userRepository.FindOneByCredential(input.Credential)
	if err != nil {
		return nil, err
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return nil, err
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * 31 * time.Hour)

	session := domain.NewSession(token, &domain.User{ID: user.ID}, expiresAt)

	newSession, err := uc.sessionRepository.AddSession(session)
	if err != nil {
		return nil, err
	}

	output := LoginOutput{
		Token:     newSession.Token,
		CreatedAt: newSession.CreatedAt,
		ExpiresAt: newSession.ExpiresAt,
	}

	return &output, nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
