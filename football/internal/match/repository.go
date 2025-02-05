package match

import "github.com/modasby/futeboxd-api/services/football/internal/match/models"

type Repository interface {
	FindOneByID(matchID int64) (*models.Match, error)
	FindBatchByID(ids []int64) ([]models.Match, error)
	GetMatchSummary(matchID int64) (*models.MatchSummary, error)
	List(where string, params []any, pageSize, pageIndex int) ([]models.Match, error)
}
