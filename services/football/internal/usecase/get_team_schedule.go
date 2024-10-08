package usecase

import (
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
)

type ScheduleOutputDTO struct {
	RequestedSeason string         `json:"requestedSeason"`
	Events          []domain.Match `json:"events"`
}

type GetTeamScheduleUseCase struct {
	espnService service.EspnService
}

func NewGetTeamScheduleUseCase(espnService service.EspnService) *GetTeamScheduleUseCase {
	return &GetTeamScheduleUseCase{espnService: espnService}
}

func (uc *GetTeamScheduleUseCase) Execute(league string, team string, season string) (*ScheduleOutputDTO, error) {

	espnSchedule, err := uc.espnService.GetTeamSchedule(league, team, season)
	if err != nil {
		return &ScheduleOutputDTO{}, err
	}

	var schedule ScheduleOutputDTO

	schedule.RequestedSeason = espnSchedule.Season.DisplayName
	schedule.Events = make([]domain.Match, 0)

	for _, espnEvent := range espnSchedule.Events {
		venue := domain.NewVenue(espnEvent.Competitions[0].Venue.FullName, espnEvent.Competitions[0].Venue.Address.City)

		var competitors []*domain.Competitor

		for _, espnCompetitor := range espnEvent.Competitions[0].Competitors {
			logo := ""

			if len(espnCompetitor.Team.Logos) > 0 {
				logo = espnCompetitor.Team.Logos[0].Href
			}

			teamID, _ := strconv.ParseInt(espnCompetitor.Team.Id, 10, 64)

			team := domain.NewTeam(
				teamID, espnCompetitor.Team.Name, espnCompetitor.Team.Abbreviation,
				espnCompetitor.Team.Color, logo,
			)

			competitor := domain.NewCompetitor(
				espnCompetitor.HomeAway, espnCompetitor.Winner, int32(espnCompetitor.Score.Value),
				team, nil,
			)

			competitors = append(competitors, competitor)
		}

		eventID, _ := strconv.ParseInt(espnEvent.Id, 10, 64)

		var note string

		if len(espnEvent.Competitions[0].Notes) > 0 {
			note = espnEvent.Competitions[0].Notes[0].Headline
		}

		newMatch := domain.NewMatch(
			eventID, *venue, 0, espnEvent.Date, note,
			competitors, espnEvent.Competitions[0].Status.Type.Completed, espnEvent.Competitions[0].Status.Type.Name,
			espnEvent.SeasonType.Name,
		)
		schedule.Events = append(schedule.Events, *newMatch)
	}

	return &schedule, nil
}
