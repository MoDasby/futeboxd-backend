package football

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/modasby/futeboxd-api/pkg/errors"
)

type Client interface {
	GetMatch(matchID int64) (*Match, error)
	GetTeam(teamID int64) (*Team, error)
}

type footballClient struct {
	baseUrl string
}

// cria um novo cliente do serviço de futebol
func NewClient(baseUrl string) Client {
	return &footballClient{
		baseUrl: baseUrl,
	}
}

func (fs *footballClient) GetMatch(matchID int64) (*Match, error) {
	res, err := http.Get(fmt.Sprintf("%s/match/%d", fs.baseUrl, matchID))
	if err != nil || res.StatusCode != 200 {
		return nil, errors.NewErrBadRequest("partida não encontrada")
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
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var team Team

	if err := json.NewDecoder(res.Body).Decode(&team); err != nil {
		log.Println(err)
		return nil, err
	}

	return &team, err
}
