package match

type Repository interface {
	Create(match *Match, summary *MatchSummary) error
	FindOneByID(matchID int64) (*Match, error)
	FindBatchByID(ids []int64) ([]Match, error)
	GetMatchSummary(matchID int64) (*MatchSummary, error)
	List(where string, params []any, pageSize, pageIndex int) ([]Match, error)
}
