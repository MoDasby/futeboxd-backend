package league

import "context"

type Repository interface {
	ListLeagues(context.Context) ([]League, error)
}
