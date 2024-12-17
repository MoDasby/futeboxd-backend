package football

type Client interface {
	GetMatch(matchID int64) (*Match, error)
	GetMatches(matchIDs []int64) ([]Match, error)
	GetMatchesMap(matchIDs []int64) (map[int64]Match, error)
	GetTeam(teamID int64) (*Team, error)
}
