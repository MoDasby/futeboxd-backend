package dto

import (
	"time"

	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
)

type ReviewInput struct {
	Rate        int    `json:"rate"`
	Description string `json:"description"`
	MatchID     int64  `json:"match_id"`
}

type ContentAuthor struct {
	ID             string      `json:"id"`
	Name           null.String `json:"name"`
	Username       string      `json:"username"`
	ProfilePicture string      `json:"profile_picture"`
}

type Review struct {
	ID            int             `json:"id"`
	Author        ContentAuthor   `json:"author"`
	Rate          int             `json:"rate"`
	Description   string          `json:"description"`
	Match         *football.Match `json:"match"`
	Likes         int             `json:"like_count"`
	CommentsCount int             `json:"comments_count"`
	IsLiked       bool            `json:"is_liked"`
	CreatedAt     time.Time       `json:"created_at"`
}
