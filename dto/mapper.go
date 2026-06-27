package dto

import (
	"spotsync/models"
)

// ToUserResponse converts a User model to UserResponse DTO.
func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

// ToZoneResponse converts a ParkingZone model to ZoneResponse DTO.
func ToZoneResponse(zone *models.ParkingZone) ZoneResponse {
	return ZoneResponse{
		ID:          zone.ID,
		Name:        zone.Name,
		Location:    zone.Location,
		Description: zone.Description,
		Capacity:    zone.Capacity,
		HourlyRate:  zone.HourlyRate,
		IsActive:    zone.IsActive,
		CreatedAt:   zone.CreatedAt,
		UpdatedAt:   zone.UpdatedAt,
	}
}

// ToZoneResponses converts a slice of ParkingZone models to ZoneResponse DTOs.
func ToZoneResponses(zones []models.ParkingZone) []ZoneResponse {
	result := make([]ZoneResponse, len(zones))
	for i, zone := range zones {
		result[i] = ToZoneResponse(&zone)
	}
	return result
}

// ToReservationResponse converts a Reservation model to ReservationResponse DTO.
func ToReservationResponse(r *models.Reservation) ReservationResponse {
	resp := ReservationResponse{
		ID:            r.ID,
		UserID:        r.UserID,
		ParkingZoneID: r.ParkingZoneID,
		VehicleNumber: r.VehicleNumber,
		StartTime:     r.StartTime,
		EndTime:       r.EndTime,
		TotalCost:     r.TotalCost,
		Status:        r.Status,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}

	if r.ParkingZone.ID != 0 {
		zone := ToZoneResponse(&r.ParkingZone)
		resp.ParkingZone = &zone
	}

	if r.User.ID != 0 {
		user := ToUserResponse(&r.User)
		resp.User = &user
	}

	return resp
}

// ToReservationResponses converts a slice of Reservation models to ReservationResponse DTOs.
func ToReservationResponses(reservations []models.Reservation) []ReservationResponse {
	result := make([]ReservationResponse, len(reservations))
	for i, r := range reservations {
		result[i] = ToReservationResponse(&r)
	}
	return result
}
