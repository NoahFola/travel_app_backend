package domain

import (
	"time"
)

type Activity struct {
	ID          string     `json:"id"`
	TripID      string     `json:"trip_id"`
	ItineraryID *string    `json:"itinerary_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Type        *string    `json:"type"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
