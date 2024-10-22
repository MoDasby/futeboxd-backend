package pagination

type Pageable[T any] struct {
	Pages       int `json:"pages,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
	Size        int `json:"size,omitempty"`
	TotalItems  int `json:"total_items,omitempty"`
	Items       []T `json:"items"`
}
