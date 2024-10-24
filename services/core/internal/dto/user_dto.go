package dto

import "github.com/modasby/futeboxd-api/pkg/client/football"

type UserDTO struct {
	ID           string         `json:"id"`
	Username     string         `json:"username"`
	Email        string         `json:"email,omitempty"`
	FavoriteTeam *football.Team `json:"favorite_team,omitempty"`
}

type UserInputDTO struct {
	Username       string `json:"username"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}
