package pagination

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type Page struct {
	Size  int
	Index int
}

func WithPage(ctx context.Context, size, index string) (*Page, error) {

	sizeInt, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		sizeInt = 50
	}

	indexInt, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		indexInt = 1
	}

	if sizeInt > 100 {
		return nil, &errors.HTTPErr{
			Msg:        "tamanho da página não pode ser maior que 100",
			Code:       400,
			Context:    "PAGINATION:PAGE_SIZE_TOO_BIG",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}
	}

	return &Page{
		Size:  int(sizeInt),
		Index: int(indexInt),
	}, nil
}

func WithRequest(r *http.Request) (*Page, error) {
	return WithPage(
		r.Context(),
		r.URL.Query().Get("page_size"),
		r.URL.Query().Get("page"),
	)
}
