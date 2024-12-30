package upload

import "context"

type Repository interface {
	Put(ctx context.Context, input *File) (string, error)
	Get(ctx context.Context, filename string) (*File, error)
}
