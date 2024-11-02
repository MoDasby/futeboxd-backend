package football

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type Client interface {
	GetMatch(matchID int64) (*Match, error)
	GetTeam(teamID int64) (*Team, error)
}

type footballClient struct {
	baseUrl string
}

func NewClient() Client {
	baseUrl := os.Getenv("FOOTBALL_URL")
	return &footballClient{
		baseUrl: baseUrl,
	}
}

func (fs *footballClient) GetMatch(matchID int64) (*Match, error) {
	res, err := http.Get(fmt.Sprintf("%s/match/%d", fs.baseUrl, matchID))
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.NewHTTPErr("partida não encontrada", 400, "FOOTBALL_CLIENT:GET_MATCH")
	}
	defer res.Body.Close()

	var match Match

	if err := json.NewDecoder(res.Body).Decode(&match); err != nil {
		return nil, err
	}

	return &match, err
}

func (fs *footballClient) GetTeam(teamID int64) (*Team, error) {
	res, err := http.Get(fmt.Sprintf("%s/team/%d", fs.baseUrl, teamID))
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.NewHTTPErr("time não encontrado", 400, "FOOTBALL_CLIENT:GET_TEAM")
	}
	defer res.Body.Close()

	var team Team

	if err := json.NewDecoder(res.Body).Decode(&team); err != nil {
		return nil, err
	}

	return &team, err
}
