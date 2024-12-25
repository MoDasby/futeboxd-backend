package pagination

import (
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/pkg/errors"
)

type Page struct {
	Size  int
	Index int
}

func NewPageWithRequest(r *http.Request) (Page, error) {
	size := r.URL.Query().Get("page_size")
	index := r.URL.Query().Get("page")

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
