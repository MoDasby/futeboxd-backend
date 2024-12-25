package match

import "github.com/modasby/futeboxd-api/services/football/pkg/errors"

type usecase struct {
	matchRepo Repository
}

type Usecase interface {
	FindOneByID(matchID int64) (*Match, error)
	FindBatchByID(ids []int64) ([]Match, error)
	GetMatchSummary(matchID int64) (*MatchSummary, error)
	List(opts ListMatchOptions, pageSize, pageIndex int) ([]Match, error)
}

func NewMatchUsecases(matchRepo Repository) Usecase {
	return &usecase{matchRepo: matchRepo}
}

func (uc *usecase) FindOneByID(matchID int64) (*Match, error) {
	return uc.matchRepo.FindOneByID(matchID)
}

func (uc *usecase) FindBatchByID(ids []int64) ([]Match, error) {
	if len(ids) > 100 {
		return nil, errors.NewHTTPErr(
			"máximo de 100 ids permitido",
			400,
			"SERVICE:MATCHES:FIND_BATCH_BY_ID:MAX_IDS_EXCEEDED",
		)
	}

	if len(ids) == 0 {
		return []Match{}, nil
	}

	return uc.matchRepo.FindBatchByID(ids)
}

func (uc *usecase) GetMatchSummary(matchID int64) (*MatchSummary, error) {
	return uc.matchRepo.GetMatchSummary(matchID)
}

func (uc *usecase) List(opts ListMatchOptions, pageSize, pageIndex int) ([]Match, error) {
	where, params := buildWhereClause(opts)

	return uc.matchRepo.List(where, params, pageSize, pageIndex)
}
