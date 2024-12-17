package domain

import "time"

type Recover struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

func (r *Recover) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now().UTC())
}
