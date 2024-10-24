package dto

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	Author    UserDTO   `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
