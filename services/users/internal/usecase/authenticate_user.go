package usecase

import (
	"time"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/users/internal/domain"
)

type AuthenticateUserUseCase struct {
	sessionRepository domain.SessionRepository
	userRepository    domain.UserRepository
}

func NewAuthenticateUserUsecase(
	sessionRepo domain.SessionRepository,
	userRepo domain.UserRepository,
) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		sessionRepository: sessionRepo,
		userRepository:    userRepo,
	}
}

type AuthenticateOutputDTO struct {
	Token     string         `json:"token"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	User      *UserOutputDTO `json:"user"`
}

func (uc *AuthenticateUserUseCase) Execute(token string) (*AuthenticateOutputDTO, error) {
	session, err := uc.sessionRepository.FindOneByToken(token)
	if err != nil {
		return nil, err
	}

	expiredAtUTC := session.ExpiresAt.UTC()
	currentTimeUTC := time.Now().UTC()

	if expiredAtUTC.Before(currentTimeUTC) {
		return nil, errors.NewErrUnauthorized("essa sessão está expirada")
	}

	output := AuthenticateOutputDTO{
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
		User: &UserOutputDTO{
			ID:           session.Owner.ID,
			Username:     session.Owner.Username,
			Email:        session.Owner.Email,
			FavoriteTeam: nil,
		},
	}

	return &output, nil
}
