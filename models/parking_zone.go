package models

import (
	"time"

	"gorm.io/gorm"
)

// ParkingZone represents the parking_zones table.
type ParkingZone struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:150;not null;index" json:"name"`
	Location    string         `gorm:"size:255;not null;index" json:"location"`
	Description string         `gorm:"size:500" json:"description"`
	Capacity    int            `gorm:"not null" json:"capacity"`
	HourlyRate  float64        `gorm:"type:decimal(10,2);not null" json:"hourly_rate"`
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name.
func (ParkingZone) TableName() string {
	return "parking_zones"
}
