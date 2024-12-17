package handler

import (
	"net/url"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/errors"
)

type Page struct {
	Size  int
	Index int
}

func NewPageWithQueryParam(params *url.Values) (Page, error) {
	size := params.Get("page_size")
	index := params.Get("page")

	sizeInt, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		sizeInt = 50
	}

	indexInt, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		indexInt = 1
	}

	if sizeInt > 100 {
		return Page{}, errors.NewHTTPErr(
			"tamanho da página não pode ser maior que 100",
			400,
			"PAGINATION:PAGE_SIZE_TOO_BIG",
		)
	}

	return Page{
		Size:  int(sizeInt),
		Index: int(indexInt),
	}, nil
}
