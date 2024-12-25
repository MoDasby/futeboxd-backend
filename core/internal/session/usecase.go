package session

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/session/dto"
)

type SessionUsecases interface {
	Login(ctx context.Context, input *dto.Login) (*dto.Session, error)
	Logout(ctx context.Context) error
}
