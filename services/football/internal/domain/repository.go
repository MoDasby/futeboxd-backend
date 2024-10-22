package domain

type MatchRepository interface {
	AddMatch(event *Match) error
	FindMatchByID(eventID string) (*Match, error)
	FindMatchBatch(ids []int64) ([]Match, error)
}

type TeamRepository interface {
	FindTeamById(ID string) (*Team, error)
	AddTeam(team *Team) error
	FindAll(name string) ([]*Team, error)
}
