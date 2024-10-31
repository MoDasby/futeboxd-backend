package dto

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	Author    UserDTO   `json:"author"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	IsLiked   bool      `json:"is_liked"`
	CreatedAt time.Time `json:"created_at"`
}
