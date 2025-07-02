package dto

import (
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
)

type ProfileStats struct {
	FollowersCount   int            `json:"followers_count"`
	FollowingCount   int            `json:"following_count"`
	Following        bool           `json:"following"`
	ReviewsCount     int            `json:"reviews_count"`
	MostReviewedTeam *football.Team `json:"most_reviewed_team"`
	AvgRating        float64        `json:"avg_rating"`
	Self             bool           `json:"self"`
}

type Profile struct {
	ID             string         `json:"id"`
	Name           null.String    `json:"name"`
	Bio            null.String    `json:"bio"`
	ProfilePicture string         `json:"profile_picture"`
	Username       string         `json:"username"`
	FavoriteTeam   *football.Team `json:"favorite_team"`
	Stats          *ProfileStats  `json:"stats"`
}
