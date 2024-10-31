package dto

import "time"

type ReviewDTO struct {
	ID          int       `json:"id"`
	Author      UserDTO   `json:"author"`
	Rate        int       `json:"rate"`
	Description string    `json:"description"`
	Match       string    `json:"match"`
	Likes       int       `json:"like_count"`
	IsLiked     bool      `json:"is_liked"`
	CreatedAt   time.Time `json:"created_at"`
}
