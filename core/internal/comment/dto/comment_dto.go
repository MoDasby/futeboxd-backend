package dto

import (
	"time"

	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
)

type CommentInput struct {
	ParentID int64  `json:"parent_id"`
	Content  string `json:"content"`
}

type ContentAuthor struct {
	ID             string      `json:"id"`
	Name           null.String `json:"name"`
	Username       string      `json:"username"`
	ProfilePicture string      `json:"profile_picture"`
}

type Comment struct {
	ID        int64         `json:"id"`
	Author    ContentAuthor `json:"author"`
	Content   string        `json:"content"`
	LikeCount int           `json:"like_count"`
	IsLiked   bool          `json:"is_liked"`
	CreatedAt time.Time     `json:"created_at"`
}
