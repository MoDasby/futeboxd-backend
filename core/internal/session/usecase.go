package session

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"

	"github.com/modasby/futeboxd-backend/core/internal/domain"
	"github.com/modasby/futeboxd-backend/core/internal/session/dto"
)

type SessionUsecases interface {
	Login(ctx context.Context, input *dto.Login) (*dto.Session, error)
	Logout(ctx context.Context) error
}

type sessionUsecases struct {
	sessionRepo    domain.SessionRepository
	userRepo       domain.UserRepository
	footballClient football.Client
}

func NewSessionUsecases(
	sessionRepo domain.SessionRepository,
	userRepo domain.UserRepository,
	footballClient football.Client,
) SessionUsecases {
	return &sessionUsecases{
		sessionRepo:    sessionRepo,
		userRepo:       userRepo,
		footballClient: footballClient,
	}
}

func (uc *sessionUsecases) Login(ctx context.Context, input *dto.Login) (*dto.Session, error) {
	user, err := uc.userRepo.FindOneByCredential(ctx, input.Credential)
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

	newSession, err := uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	output := dto.Session{
		Token:     newSession.Token,
		CreatedAt: newSession.CreatedAt,
		ExpiresAt: newSession.ExpiresAt,
	}

	return &output, nil
}

func (uc *sessionUsecases) Logout(ctx context.Context) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if err := uc.sessionRepo.Delete(ctx, session.ID); err != nil {
		return err
	}

	return nil
}
