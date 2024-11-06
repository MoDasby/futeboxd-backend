package domain

type MatchRepository interface {
	Create(event *Match) error
	FindOneByID(eventID string) (*Match, error)
	FindBatchByID(ids []int64) ([]Match, error)
}

type TeamRepository interface {
	FindOneById(ID string) (*Team, error)
	Create(team *Team) error
	FindAll(name string) ([]*Team, error)
}
