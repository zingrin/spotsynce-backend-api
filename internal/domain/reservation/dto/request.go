package dto

import "github.com/google/uuid"

type CreateRequest struct {
	ZoneId       uuid.UUID `json:"zone_id" validate:"required"`
	LicensePlate string    `json:"license_plate" validate:"required,max=15"`
}
