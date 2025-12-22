package domain

import (
	"time"
)

type Itinerary struct {
	ID        string    `json:"id"`
	TripID    string    `json:"trip_id"`
	Slug      string    `json:"slug"`
	Title     *string   `json:"title"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
