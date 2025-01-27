package upload

import (
	"context"
	"io"
)

type Usecase interface {
	Upload(ctx context.Context, file io.ReadSeeker) error
}
