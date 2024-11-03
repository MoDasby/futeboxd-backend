package pagination

import (
	"strconv"

	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type Page struct {
	Size  int
	Index int
}

func WithPage(size, index string) (*Page, error) {

	sizeInt, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		sizeInt = 50
	}

	indexInt, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		indexInt = 1
	}

	if sizeInt > 100 {
		return nil, errors.NewHTTPErr(
			"tamanho da página não pode ser maior que 100",
			400,
			"PAGINATION:PAGE_SIZE_TOO_BIG",
		)
	}

	return &Page{
		Size:  int(sizeInt),
		Index: int(indexInt),
	}, nil
}
