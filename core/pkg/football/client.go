package football

import "context"

type Client interface {
	GetMatch(ctx context.Context, matchID int64) (*Match, error)
	GetMatches(ctx context.Context, matchIDs []int64) ([]Match, error)
	GetMatchesMap(ctx context.Context, matchIDs []int64) (map[int64]Match, error)
	GetTeam(ctx context.Context, teamID int64) (*Team, error)
}
