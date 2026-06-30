package dto

import (
	"time"

	"github.com/google/uuid"
)

type ZoneResponse struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	TotalCapacity int       `json:"total_capacity"`
	PricePerHour  float64   `json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ZoneWithAvailability struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	TotalCapacity  int       `json:"total_capacity"`
	AvailableSpots int       `json:"available_spots"`
	PricePerHour   float64   `json:"price_per_hour"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
