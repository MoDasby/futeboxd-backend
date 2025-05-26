package sq

import "github.com/Masterminds/squirrel"

type NotEq = squirrel.NotEq
type Eq = squirrel.Eq
type ILike = squirrel.ILike

var Query = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
