package dto

import (
	"time"

	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
)

type ReviewDTO struct {
	ID            int             `json:"id"`
	Author        UserDTO         `json:"author"`
	Rate          int             `json:"rate"`
	Description   string          `json:"description"`
	Match         *football.Match `json:"match"`
	Likes         int             `json:"like_count"`
	CommentsCount int             `json:"comments_count"`
	IsLiked       bool            `json:"is_liked"`
	CreatedAt     time.Time       `json:"created_at"`
}
