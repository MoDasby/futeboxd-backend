package dto

import (
	"github.com/modasby/futeboxd-api/services/core/internal/json/null"
)

type UserDTO struct {
	ID       string      `json:"id"`
	Name     null.String `json:"name"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
}

type UserInputDTO struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	Bio            string `json:"bio"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}
