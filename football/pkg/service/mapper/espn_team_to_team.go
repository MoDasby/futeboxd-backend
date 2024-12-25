package mapper

import (
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/team"
	"github.com/modasby/futeboxd-api/services/football/pkg/service/espn"
)

func EspnTeamToTeam(espnTeam espn.EspnTeam) (*team.Team, error) {
	ID, err := strconv.ParseUint(espnTeam.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	var logo string

	if len(espnTeam.Logos) > 0 {
		logo = espnTeam.Logos[0].Href
	}

	return &team.Team{
		ID:           int64(ID),
		Name:         espnTeam.Name,
		Abbreviation: espnTeam.Abbreviation,
		Color:        espnTeam.Color,
		Logo:         logo,
	}, nil
}
