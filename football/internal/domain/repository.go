package domain

type MatchRepository interface {
	Create(match *Match, summary *MatchSummary) error
	FindOneByID(matchID int64) (*Match, error)
	FindBatchByID(ids []int64) ([]Match, error)
	GetMatchSummary(matchID int64) (*MatchSummary, error)
	ListByYear(teamID int64, year, pageSize, pageIndex int) ([]Match, error)
	FindLiveMatches(pageSize, pageIndex int) ([]Match, error)
}

type TeamRepository interface {
	FindOneById(teamID int64) (*Team, error)
	Create(team *Team) error
	FindAll(name string, pageSize, pageIndex int) ([]Team, error)
}
