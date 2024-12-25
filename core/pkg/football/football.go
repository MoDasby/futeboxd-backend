package football

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-backend/core/pkg/errors"
)

type footballClient struct {
	baseUrl string
}

func NewClient() Client {
	baseUrl := os.Getenv("FOOTBALL_URL")
	return &footballClient{
		baseUrl: baseUrl,
	}
}

func (fs *footballClient) GetMatchesMap(matchIDs []int64) (map[int64]Match, error) {
	matches, err := fs.GetMatches(matchIDs)
	if err != nil {
		return nil, err
	}

	matchesMap := make(map[int64]Match)

	for _, match := range matches {
		matchesMap[match.ID] = match
	}

	return matchesMap, nil
}

func (fs *footballClient) GetMatches(matchIDs []int64) ([]Match, error) {
	body, err := json.Marshal(matchIDs)
	if err != nil {
		return nil, err
	}

	// TODO setar o header User-Agent pra core client

	res, err := http.Post(fmt.Sprintf("%s/matches/batch", fs.baseUrl), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body errors.ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, errors.NewHTTPErr(body.Msg, res.StatusCode, "CLIENT:FOOTBALL:GET_MATCHES")
	}

	matches := make([]Match, 0)

	if err := json.NewDecoder(res.Body).Decode(&matches); err != nil {
		return nil, err
	}

	return matches, err
}

func (fs *footballClient) GetMatch(matchID int64) (*Match, error) {
	res, err := http.Get(fmt.Sprintf("%s/matches/%d", fs.baseUrl, matchID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body errors.ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, errors.NewHTTPErr(body.Msg, res.StatusCode, "CLIENT:FOOTBALL:GET_MATCH")
	}

	var match Match

	if err := json.NewDecoder(res.Body).Decode(&match); err != nil {
		return nil, err
	}

	return &match, err
}

func (fs *footballClient) GetTeam(teamID int64) (*Team, error) {
	res, err := http.Get(fmt.Sprintf("%s/teams/%d", fs.baseUrl, teamID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body errors.ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, errors.NewHTTPErr(body.Msg, res.StatusCode, "CLIENT:FOOTBALL:GET_TEAM")
	}

	var team Team

	if err := json.NewDecoder(res.Body).Decode(&team); err != nil {
		return nil, err
	}

	return &team, err
}
