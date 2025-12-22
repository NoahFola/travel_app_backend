package domain

import (
	"time"
)

type Trip struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Location  string    `json:"location"`
	StartDate time.Time `json:"start_date"` // Keeping as time.Time, usually handled as date in logic
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
