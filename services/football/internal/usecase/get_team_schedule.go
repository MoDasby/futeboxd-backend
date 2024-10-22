package usecase

import (
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
)

type ScheduleOutputDTO struct {
	RequestedSeason string         `json:"requestedSeason"`
	Matches         []domain.Match `json:"matches"`
}

type GetTeamScheduleUseCase struct {
	espnService service.EspnService
}

func NewGetTeamScheduleUseCase(espnService service.EspnService) *GetTeamScheduleUseCase {
	return &GetTeamScheduleUseCase{espnService: espnService}
}

func (uc *GetTeamScheduleUseCase) Execute(league, team, season string) (*ScheduleOutputDTO, error) {

	espnSchedule, err := uc.espnService.GetTeamSchedule(league, team, season)
	if err != nil {
		return &ScheduleOutputDTO{}, err
	}

	var schedule ScheduleOutputDTO

	schedule.RequestedSeason = espnSchedule.Season.DisplayName
	schedule.Matches = make([]domain.Match, 0)

	for _, espnEvent := range espnSchedule.Events {
		venue := domain.NewVenue(espnEvent.Competitions[0].Venue.FullName, espnEvent.Competitions[0].Venue.Address.City)

		var homeCompetitor domain.Competitor
		var awayCompetitor domain.Competitor

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
				*team, nil,
			)

			if competitor.HomeAway == "home" {
				homeCompetitor = *competitor
			}

			if competitor.HomeAway == "away" {
				awayCompetitor = *competitor
			}
		}

		eventID, _ := strconv.ParseInt(espnEvent.Id, 10, 64)

		var note string

		if len(espnEvent.Competitions[0].Notes) > 0 {
			note = espnEvent.Competitions[0].Notes[0].Headline
		}

		newMatch := domain.NewMatch(
			eventID, *venue, espnEvent.Date, note,
			homeCompetitor, awayCompetitor, espnEvent.Competitions[0].Status.Type.Completed, espnEvent.Competitions[0].Status.Type.Name,
			espnEvent.SeasonType.Name, nil,
		)
		schedule.Matches = append(schedule.Matches, *newMatch)
	}

	return &schedule, nil
}
