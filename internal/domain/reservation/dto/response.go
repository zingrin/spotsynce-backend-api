package dto

import (
	"time"

	"github.com/google/uuid"
)

type ReservationResponse struct {
	Id           uuid.UUID `json:"id"`
	UserId       uuid.UUID `json:"user_id"`
	ZoneId       uuid.UUID `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserInfo struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type ZoneInfo struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

type MyReservationResponse struct {
	Id           uuid.UUID `json:"id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	Zone         ZoneInfo  `json:"zone"`
	CreatedAt    time.Time `json:"created_at"`
}

type AdminReservationResponse struct {
	Id           uuid.UUID `json:"id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	User         UserInfo  `json:"user"`
	Zone         ZoneInfo  `json:"zone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
