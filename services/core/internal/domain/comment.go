package domain

import "time"

type Comment struct {
	ID        int64
	Author    User
	ParentID  int64
	Content   string
	CreatedAt time.Time
}

func NewComment(
	author User,
	parentID int64,
	content string,
) *Comment {
	return &Comment{
		Author:   author,
		ParentID: parentID,
		Content:  content,
	}
}
