package domain

import "time"

type Comment struct {
	ID        int64
	Author    User
	Comments  []Comment
	Content   string
	CreatedAt time.Time
}
