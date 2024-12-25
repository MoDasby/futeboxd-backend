package user

import "context"

type Repository interface {
	FindOneByCredential(ctx context.Context, credential string) (*User, error)
	Create(ctx context.Context, user *User) error
	Exists(ctx context.Context, username, email string) (bool, error)
	FindOneByIdOrUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	SaveRecoverToken(ctx context.Context, recover *Recover) error
	CheckRecoverToken(ctx context.Context, token string) (*Recover, error)
}
