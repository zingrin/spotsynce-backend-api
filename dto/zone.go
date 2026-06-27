package dto

import "time"

// CreateZoneRequest holds parking zone creation input.
type CreateZoneRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=150"`
	Location    string  `json:"location" validate:"required,min=2,max=255"`
	Description string  `json:"description" validate:"omitempty,max=500"`
	Capacity    int     `json:"capacity" validate:"required,gt=0"`
	HourlyRate  float64 `json:"hourly_rate" validate:"required,gte=0"`
}

// ZoneResponse holds parking zone output data.
type ZoneResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Capacity    int       `json:"capacity"`
	HourlyRate  float64   `json:"hourly_rate"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ZoneListQuery holds query parameters for listing parking zones.
type ZoneListQuery struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Search   string `query:"search"`
	Location string `query:"location"`
	SortBy   string `query:"sort_by"`
	SortDir  string `query:"sort_dir"`
	IsActive *bool  `query:"is_active"`
}

// PaginatedResponse wraps paginated list results.
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
