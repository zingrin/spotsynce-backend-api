package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ReservationStatusActive    = "active"
	ReservationStatusCancelled = "cancelled"
	ReservationStatusCompleted = "completed"
)

// Reservation represents the reservations table.
type Reservation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	ParkingZoneID uint           `gorm:"not null;index" json:"parking_zone_id"`
	VehicleNumber string         `gorm:"size:20;not null" json:"vehicle_number"`
	StartTime     time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime       time.Time      `gorm:"not null;index" json:"end_time"`
	TotalCost     float64        `gorm:"type:decimal(10,2);not null" json:"total_cost"`
	Status        string         `gorm:"size:20;not null;default:active;index" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParkingZone ParkingZone `gorm:"foreignKey:ParkingZoneID" json:"parking_zone,omitempty"`
}

// TableName overrides the default table name.
func (Reservation) TableName() string {
	return "reservations"
}
