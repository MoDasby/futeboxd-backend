package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/internal/session/dto"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type sessionUsecases struct {
	sessionRepo session.Repository
	userRepo    user.Repository
}

func NewSessionUsecases(
	sessionRepo session.Repository,
	userRepo user.Repository,
) session.SessionUsecases {
	return &sessionUsecases{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

func (uc *sessionUsecases) Login(ctx context.Context, input *dto.Login) (*dto.Session, error) {
	user, err := uc.userRepo.FindOneByCredential(ctx, input.Credential)
	if err != nil {
		return nil, &errors.HTTPErr{
			Msg:        "Usuário ou senha inválidos",
			Code:       http.StatusBadRequest,
			Context:    "SESSION:USECASE:LOGIN:USER_NOT_FOUND",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}
	}

	if err := user.CheckPassword(ctx, input.Password); err != nil {
		return nil, &errors.HTTPErr{
			Msg:        "Usuário ou senha inválidos",
			Code:       http.StatusBadRequest,
			Context:    "SESSION:USECASE:LOGIN:WRONG_PASSWORD",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}
	}

	session, err := session.NewSession(user.ID)
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
