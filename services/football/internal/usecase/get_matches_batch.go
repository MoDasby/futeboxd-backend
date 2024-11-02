package usecase

import (
	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/errors"
)

type GetMatchesBatchUseCase struct {
	matchRepository domain.MatchRepository
}

func NewGetMatchesBatchUseCase(repo domain.MatchRepository) *GetMatchesBatchUseCase {
	return &GetMatchesBatchUseCase{
		matchRepository: repo,
	}
}

type GetMatchesBatchInput struct {
	IDs []int64 `json:"ids"`
}

func (uc GetMatchesBatchUseCase) Execute(input GetMatchesBatchInput) ([]domain.Match, error) {
	if len(input.IDs) == 0 {
		return nil, errors.NewHTTPErr("corpo de requisição inválido", 400, "USECASE:GET_MATCH_BATCH:INVALID_BODY")
	}

	matches, err := uc.matchRepository.FindMatchBatch(input.IDs)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
