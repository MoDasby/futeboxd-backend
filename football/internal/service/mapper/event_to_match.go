package mapper

import (
	"log"
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/service/espn"
)

func EspnEventToMatch(espnEvent *espn.EspnEventSummary) (*domain.Match, *domain.MatchSummary, error) {

	var homeCompetitor domain.Competitor
	var awayCompetitor domain.Competitor

	var summary domain.MatchSummary

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
			return nil, nil, err
		}

		espnRoster := espnEvent.Rosters[index].Roster

		roster := make([]domain.Roster, 0)

		for _, r := range espnRoster {
			newRoster := domain.NewRoster(r.Starter, r.Athlete.ID, r.Athlete.LastName, r.Athlete.FullName, r.Athlete.DisplayName, r.Athlete.HeadShot.Href, r.Athlete.HeadShot.Alt, r.Position.DisplayName, r.Position.Abbreviation, r.SubbedIn, r.SubbedOut)
			roster = append(roster, *newRoster)
		}

		competitor := domain.NewCompetitor(
			espnCompetitor.HomeAway, espnCompetitor.Winner, int32(score),
			*team,
		)

		if competitor.HomeAway == "home" {
			homeCompetitor = *competitor
			summary.HomeCompetitorRoster = roster
		}

		if competitor.HomeAway == "away" {
			awayCompetitor = *competitor
			summary.AwayCompetitorRoster = roster
		}
	}

	matchID, _ := strconv.ParseInt(espnEvent.Header.ID, 10, 64)

	var note string

	if len(espnEvent.Header.Competitions[0].Notes) > 0 {
		note = espnEvent.Header.Competitions[0].Notes[0].Headline
	}

	var events []domain.Event

	for _, keyEvent := range espnEvent.KeyEvents {

		clockValue, err := strconv.ParseInt(keyEvent.Clock.DisplayValue, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		events = append(events, *domain.NewEvent(
			keyEvent.Type.ID,
			keyEvent.Type.Text,
			keyEvent.Text,
			int(clockValue),
		))
	}

	summary.Events = events

	match := domain.NewMatch(
		matchID, espnEvent.GameInfo.Venue.FullName, espnEvent.Header.Competitions[0].Date,
		note, homeCompetitor, awayCompetitor, espnEvent.Header.Competitions[0].Status.Type.Completed,
		espnEvent.Header.Competitions[0].Status.Type.Name,
		espnEvent.Header.Season.Name,
	)

	return match, &summary, nil
}
