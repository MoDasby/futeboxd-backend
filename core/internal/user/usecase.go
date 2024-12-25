package user

import (
	"context"

	"github.com/modasby/futeboxd-backend/core/internal/user/dto"
)

type Usecases interface {
	CreateUser(ctx context.Context, input *dto.UserInput) error
	EditUser(ctx context.Context, input *dto.EditUser) error
	ChangePassword(ctx context.Context, input *dto.ChangePassword) error
	RecoverPassword(ctx context.Context, input *dto.RecoverPassword) error
	ResetPassword(ctx context.Context, input *dto.ResetPassword) error
	GetCurrentUser(ctx context.Context) (*dto.User, error)
}
