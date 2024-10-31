package usecase

import (
	"time"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
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
		return nil, errors.NewHTTPErr("usuário ou senha inválidos", 400, "LOGIN:USER_NOT_FOUND")
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return nil, errors.NewHTTPErr("usuário ou senha inválidos", 400, "LOGIN:WRONG_PASSWORD")
	}

	session, err := domain.NewSession(user.ID)
	if err != nil {
		return nil, err
	}

	newSession, err := uc.sessionRepository.Create(session)
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
