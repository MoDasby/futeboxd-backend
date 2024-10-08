package user

import (
	"time"

	"github.com/modasby/futeboxd-api/pkg/client/football"
)

type UserOutputDTO struct {
	ID           string         `json:"id"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	FavoriteTeam *football.Team `json:"favorite_team,omitempty"`
}

type AuthenticateOutputDTO struct {
	Token     string         `json:"token"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	User      *UserOutputDTO `json:"user"`
}
