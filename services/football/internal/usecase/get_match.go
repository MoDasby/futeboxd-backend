package usecase

import (
	"log"
	"strconv"

	"github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
)

type GetMatchUseCase struct {
	espnService     service.EspnService
	matchRepository domain.MatchRepository
}

func NewGetMatchUseCase(espnService service.EspnService, matchRepository domain.MatchRepository) *GetMatchUseCase {
	return &GetMatchUseCase{
		espnService:     espnService,
		matchRepository: matchRepository,
	}
}

func (uc *GetMatchUseCase) Execute(id string) (*domain.Match, error) {
	existingEvent, err := uc.matchRepository.FindMatchByID(id)
	if err != nil {
		// se a partida não for encontrada não retorna um erro
		if _, ok := err.(*errors.ErrNotFound); !ok {
			return nil, err
		}
	}

	if existingEvent != nil {
		return existingEvent, nil
	}

	espnEvent, err := uc.espnService.GetEvent(id)
	if err != nil {
		return nil, err
	}

	venue := domain.NewVenue(espnEvent.GameInfo.Venue.FullName, espnEvent.GameInfo.Venue.Address.City)

	var competitors []*domain.Competitor

	for index, espnCompetitor := range espnEvent.Header.Competitions[0].Competitors {
		logo := ""

		if len(espnCompetitor.Team.Logos) > 0 {
			logo = espnCompetitor.Team.Logos[0].Href
		}

		teamID, _ := strconv.ParseInt(espnCompetitor.Team.Id, 10, 64)

		team := domain.NewTeam(
			teamID, espnCompetitor.Team.Name, espnCompetitor.Team.Abbreviation,
			espnCompetitor.Team.Color, logo,
		)

		score, err := strconv.Atoi(espnCompetitor.Score)
		if err != nil {
			return nil, err
		}

		espnRoster := espnEvent.Rosters[index].Roster

		roster := make([]*domain.Roster, 0)

		for _, r := range espnRoster {
			newRoster := domain.NewRoster(r.Starter, r.Athlete.ID, r.Athlete.LastName, r.Athlete.FullName, r.Athlete.DisplayName, r.Athlete.HeadShot.Href, r.Athlete.HeadShot.Alt, r.Position.DisplayName, r.Position.Abbreviation, r.SubbedIn, r.SubbedOut)
			roster = append(roster, newRoster)
		}

		competitor := domain.NewCompetitor(
			espnCompetitor.HomeAway, espnCompetitor.Winner, int32(score),
			team, roster,
		)

		competitors = append(competitors, competitor)
	}

	matchID, _ := strconv.ParseInt(espnEvent.Header.ID, 10, 64)

	var note string

	if len(espnEvent.Header.Competitions[0].Notes) > 0 {
		note = espnEvent.Header.Competitions[0].Notes[0].Headline
	}

	match := domain.NewMatch(
		matchID, *venue, 0, espnEvent.Header.Competitions[0].Date,
		note, competitors, espnEvent.Header.Competitions[0].Status.Type.Completed,
		espnEvent.Header.Competitions[0].Status.Type.Name,
		espnEvent.Header.Season.Name,
	)

	if err := uc.matchRepository.AddMatch(match); err != nil {
		log.Print(err)
	}

	return match, nil
}
