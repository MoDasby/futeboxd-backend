package dto

import "github.com/modasby/futeboxd-api/services/core/internal/client/football"

type UserDTO struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Bio          string         `json:"bio"`
	Username     string         `json:"username"`
	Email        string         `json:"email,omitempty"`
	FavoriteTeam *football.Team `json:"favorite_team,omitempty"`
}

type UserInputDTO struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	Bio            string `json:"bio"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}
