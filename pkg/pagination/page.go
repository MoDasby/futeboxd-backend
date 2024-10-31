package pagination

import (
	"strconv"
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

	return &Page{
		Size:  int(sizeInt),
		Index: int(indexInt),
	}, nil
}
