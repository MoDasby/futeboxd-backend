package dto

import "time"

type CommentInput struct {
	ParentID int64  `json:"parent_id"`
	Content  string `json:"content"`
}

type ContentAuthor struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Comment struct {
	ID        int64         `json:"id"`
	Author    ContentAuthor `json:"author"`
	Content   string        `json:"content"`
	LikeCount int           `json:"like_count"`
	IsLiked   bool          `json:"is_liked"`
	CreatedAt time.Time     `json:"created_at"`
}
