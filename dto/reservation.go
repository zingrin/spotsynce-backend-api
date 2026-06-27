package dto

import "time"

// CreateReservationRequest holds reservation creation input.
type CreateReservationRequest struct {
	ParkingZoneID uint      `json:"parking_zone_id" validate:"required,gt=0"`
	VehicleNumber string    `json:"vehicle_number" validate:"required,min=2,max=20"`
	StartTime     time.Time `json:"start_time" validate:"required"`
	EndTime       time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
}

// ReservationResponse holds reservation output data.
type ReservationResponse struct {
	ID            uint         `json:"id"`
	UserID        uint         `json:"user_id"`
	ParkingZoneID uint         `json:"parking_zone_id"`
	VehicleNumber string       `json:"vehicle_number"`
	StartTime     time.Time    `json:"start_time"`
	EndTime       time.Time    `json:"end_time"`
	TotalCost     float64      `json:"total_cost"`
	Status        string       `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	ParkingZone   *ZoneResponse `json:"parking_zone,omitempty"`
	User          *UserResponse `json:"user,omitempty"`
}

// ReservationListQuery holds query parameters for listing reservations.
type ReservationListQuery struct {
	Page          int    `query:"page"`
	Limit         int    `query:"limit"`
	Status        string `query:"status"`
	ParkingZoneID uint   `query:"parking_zone_id"`
	UserID        uint   `query:"user_id"`
	SortBy        string `query:"sort_by"`
	SortDir       string `query:"sort_dir"`
}
