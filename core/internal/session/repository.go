package session

import "context"

type Repository interface {
	FindOneByToken(ctx context.Context, token string) (*Session, error)
	Create(ctx context.Context, session *Session) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
}
