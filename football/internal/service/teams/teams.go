package teams

import (
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/service/espn"
)

type teamService struct {
	teamRepo    domain.TeamRepository
	espnService espn.EspnService
}

type TeamService interface {
	FindOneById(teamID int64) (*domain.Team, error)
	Create(team *domain.Team) error
	FindAll(name string, pageSize, pageIndex int) ([]domain.Team, error)
	ListSchedule(team, season string) (*ScheduleOutputDTO, error)
}

func NewTeamService(teamRepo domain.TeamRepository, espnService espn.EspnService) TeamService {
	return &teamService{
		teamRepo:    teamRepo,
		espnService: espnService,
	}
}

func (s *teamService) FindOneById(teamID int64) (*domain.Team, error) {
	return s.teamRepo.FindOneById(teamID)
}

func (s *teamService) Create(team *domain.Team) error {
	return s.teamRepo.Create(team)
}

func (s *teamService) FindAll(name string, pageSize, pageIndex int) ([]domain.Team, error) {
	return s.teamRepo.FindAll(name, pageSize, pageIndex)
}

type ScheduleOutputDTO struct {
	RequestedSeason string         `json:"requestedSeason"`
	Matches         []domain.Match `json:"matches"`
}

func (s *teamService) ListSchedule(team, season string) (*ScheduleOutputDTO, error) {
	espnSchedule, err := s.espnService.GetTeamSchedule("all", team, season)
	if err != nil {
		return &ScheduleOutputDTO{}, err
	}

	var schedule ScheduleOutputDTO

	schedule.RequestedSeason = espnSchedule.Season.DisplayName
	schedule.Matches = make([]domain.Match, 0)

	for _, espnEvent := range espnSchedule.Events {

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
				*team,
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
			eventID, espnEvent.Competitions[0].Venue.FullName, espnEvent.Date, note,
			homeCompetitor, awayCompetitor, espnEvent.Competitions[0].Status.Type.Completed, espnEvent.Competitions[0].Status.Type.Name,
			espnEvent.SeasonType.Name,
		)
		schedule.Matches = append(schedule.Matches, *newMatch)
	}

	return &schedule, nil
}
