package dto

import (
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
)

type User struct {
	ID             string         `json:"id"`
	Name           null.String    `json:"name"`
	Username       string         `json:"username"`
	Bio            null.String    `json:"bio"`
	Email          string         `json:"email"`
	FavoriteTeam   *football.Team `json:"favorite_team"`
	ProfilePicture string         `json:"profile_picture"`
}

type UserInput struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	Bio            string `json:"bio"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}

type EditUser struct {
	Name           string `json:"name"`
	Bio            string `json:"bio"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FavoriteTeamID int64  `json:"favorite_team_id"`
}

type ChangePassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type RecoverPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
