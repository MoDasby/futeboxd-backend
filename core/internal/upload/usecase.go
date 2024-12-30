package upload

import (
	"context"
)

type Usecase interface {
	Upload(ctx context.Context, img File) error
}
