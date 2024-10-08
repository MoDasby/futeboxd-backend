package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-api/pkg/errors"
)

const (
	BASE_URL = "https://site.api.espn.com/apis/site/v2"
)

type EspnService interface {
	ListTeams(league string) (*EspnTeams, error)
	ListLeagues() (*EspnLeagues, error)
	GetTeamSchedule(league string, team string, season string) (*EspnSchedule, error)
	GetEvent(ID string) (*EspnEventSummary, error)
}

type espnService struct {
}

func NewEspnService() EspnService {
	return &espnService{}
}

func getResource[T any](url string) (*T, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.NewErrBadRequest("erro desconhecido ao buscar dados")
	}

	var response T

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *espnService) ListTeams(league string) (*EspnTeams, error) {
	path := fmt.Sprintf("%s/sports/soccer/%s/teams?lang=pt", BASE_URL, league)

	teams, err := getResource[EspnTeams](path)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (s *espnService) ListLeagues() (*EspnLeagues, error) {
	url := fmt.Sprintf("%s/leagues/dropdown?lang=pt&region=pt&calendartype=whitelist&limit=100&sport=soccer", BASE_URL)

	leagues, err := getResource[EspnLeagues](url)
	if err != nil {
		return nil, err
	}

	return leagues, nil
}

func (s *espnService) GetTeamSchedule(league string, team string, season string) (*EspnSchedule, error) {
	if len(season) == 0 {
		season = fmt.Sprintf("%d", time.Now().Year())
	}

	url := fmt.Sprintf("%s/sports/soccer/%s/teams/%s/schedule?lang=pt&season=%s", BASE_URL, league, team, season)
	espnSchedule, err := getResource[EspnSchedule](url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return espnSchedule, nil
}

func (s *espnService) GetEvent(ID string) (*EspnEventSummary, error) {
	url := fmt.Sprintf("%s/sports/soccer/all/summary?lang=pt&event=%s", BASE_URL, ID)

	event, err := getResource[EspnEventSummary](url)
	if err != nil {
		return nil, err
	}

	return event, nil
}
