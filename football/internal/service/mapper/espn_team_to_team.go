package mapper

import (
	"strconv"

	"github.com/modasby/futeboxd-api/services/football/internal/domain"
	"github.com/modasby/futeboxd-api/services/football/internal/service"
)

func EspnTeamToTeam(espnTeam service.EspnTeam) (*domain.Team, error) {
	ID, err := strconv.ParseUint(espnTeam.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	var logo string

	if len(espnTeam.Logos) > 0 {
		logo = espnTeam.Logos[0].Href
	}

	return &domain.Team{
		ID:           int64(ID),
		Name:         espnTeam.Name,
		Abbreviation: espnTeam.Abbreviation,
		Color:        espnTeam.Color,
		Logo:         logo,
	}, nil
}
