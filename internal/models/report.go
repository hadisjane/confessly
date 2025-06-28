package models

import "time"

type Report struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ConfessionID int     `json:"confession_id" db:"confession_id"`
	Reason    string    `json:"reason" db:"reason"`
	Status    string    `json:"status" db:"status"` // "pending", "approved", "rejected"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
}

type UpdateReport struct {
	Status *string `json:"status" db:"status"`
}