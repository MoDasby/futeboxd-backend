package pagination

import "github.com/modasby/futeboxd-api/pkg/errors"

type Page struct {
	Size  int
	Index int
}

func WithPage(size, index int) (*Page, error) {
	if size > 100 {
		return nil, errors.NewHTTPErr(
			"cada página deve ter um valor máximo de 100",
			400,
			"PAGINATION:INVALID_PAGE_SIZE",
		)
	}

	return &Page{
		Size:  size,
		Index: index,
	}, nil
}
