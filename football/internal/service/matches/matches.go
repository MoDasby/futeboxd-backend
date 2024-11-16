package matches

import (
	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/errors"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
	"github.com/modasby/futeboxd-api/services/football/internal/service/mapper"
)

type matchService struct {
	matchRepo   domain.MatchRepository
	espnService service.EspnService
}

type MatchService interface {
	FindOneByID(matchID int64) (*domain.Match, error)
	FindBatchByID(ids []int64) ([]domain.Match, error)
	GetMatchSummary(matchID int64) (*domain.MatchSummary, error)
}

func NewMatchService(
	matchRepo domain.MatchRepository,
	espnService service.EspnService,
) MatchService {
	return &matchService{matchRepo: matchRepo, espnService: espnService}
}

func (s *matchService) FindOneByID(matchID int64) (*domain.Match, error) {
	existingMatch, err := s.matchRepo.FindOneByID(matchID)
	if err != nil {
		if e, ok := err.(*errors.HTTPErr); !ok || e.Code != 404 {
			return nil, err
		}
	}

	if existingMatch != nil {
		return existingMatch, nil
	}

	espnEvent, err := s.espnService.GetEvent(matchID)
	if err != nil {
		return nil, err
	}

	match, summary, err := mapper.EspnEventToMatch(espnEvent)
	if err != nil {
		return nil, err
	}

	if err := s.matchRepo.Create(match, summary); err != nil {
		return nil, err
	}

	return match, nil
}

func (s *matchService) FindBatchByID(ids []int64) ([]domain.Match, error) {
	if len(ids) > 100 {
		return nil, errors.NewHTTPErr(
			"máximo de 100 ids permitido",
			400,
			"SERVICE:MATCHES:FIND_BATCH_BY_ID:MAX_IDS_EXCEEDED",
		)
	}

	if len(ids) == 0 {
		return []domain.Match{}, nil
	}

	return s.matchRepo.FindBatchByID(ids)
}

func (s *matchService) GetMatchSummary(matchID int64) (*domain.MatchSummary, error) {
	return s.matchRepo.GetMatchSummary(matchID)
}
