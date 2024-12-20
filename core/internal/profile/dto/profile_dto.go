package dto

import (
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
)

type Profile struct {
	ID             string         `json:"id"`
	Name           null.String    `json:"name"`
	Bio            null.String    `json:"bio"`
	Username       string         `json:"username"`
	FavoriteTeam   *football.Team `json:"favorite_team"`
	FollowersCount int            `json:"followers_count"`
	FollowingCount int            `json:"following_count"`
	Following      bool           `json:"following"`
	Self           bool           `json:"self"`
}
