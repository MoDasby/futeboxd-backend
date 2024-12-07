package dto

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/json/null"
)

type UserDTO struct {
	ID           string         `json:"id"`
	Name         null.String    `json:"name"`
	Bio          null.String    `json:"bio"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	FavoriteTeam *football.Team `json:"favorite_team"`
}

type UserInputDTO struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	Bio            string `json:"bio"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}
