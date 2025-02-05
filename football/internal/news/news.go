package news

import "time"

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	ImageLink   string    `json:"image_link"`
	CreatedAt   time.Time `json:"created_at"`
}
