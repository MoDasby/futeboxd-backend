package match

import "github.com/modasby/futeboxd-api/services/football/internal/match/models"

type ListMatchOptions struct {
	Team   int64
	Year   int
	Status string
	League int
}

type Repository interface {
	FindOneByID(matchID int64) (*models.Match, error)
	FindBatchByID(ids []int64) ([]models.Match, error)
	GetMatchSummary(matchID int64) (*models.MatchSummary, error)
	List(filters ListMatchOptions, pageSize, pageIndex int) ([]models.Match, error)
}
