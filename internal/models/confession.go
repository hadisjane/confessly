package models

import "time"

type Confession struct {
	ID        int        `json:"id" db:"id"`
	UserID    *int       `json:"user_id,omitempty" db:"user_id"`
	GuestUUID *string    `json:"guest_uuid,omitempty" db:"guest_uuid"`
	Username  string     `json:"username,omitempty" db:"username"`
	Title     string     `json:"title" binding:"required,min=5,max=100" db:"title"`
	Text      string     `json:"text" binding:"required" db:"text"`
	Anon      bool       `json:"anon" db:"anon"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}