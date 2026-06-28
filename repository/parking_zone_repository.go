package repository

import (
	"context"
	"errors"
	"strings"

	"spotsync/dto"
	"spotsync/models"

	"gorm.io/gorm"
)

// ZoneFilter holds filter criteria for parking zone queries.
type ZoneFilter struct {
	Search   string
	Location string
	IsActive *bool
	SortBy   string
	SortDir  string
	Page     int
	Limit    int
}

// ParkingZoneRepository handles all database operations for parking zones.
type ParkingZoneRepository struct {
	db *gorm.DB
}

// NewParkingZoneRepository creates a new ParkingZoneRepository.
func NewParkingZoneRepository(db *gorm.DB) *ParkingZoneRepository {
	return &ParkingZoneRepository{db: db}
}

// Create inserts a new parking zone record.
func (r *ParkingZoneRepository) Create(ctx context.Context, zone *models.ParkingZone) error {
	return r.db.WithContext(ctx).Create(zone).Error
}

// FindByID retrieves a parking zone by primary key
func (r *ParkingZoneRepository) FindByID(ctx context.Context, id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.WithContext(ctx).First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &zone, nil
}

// List retrieves parking zones with pagination, filtering, sorting, and search.
func (r *ParkingZoneRepository) List(ctx context.Context, filter ZoneFilter) ([]models.ParkingZone, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.ParkingZone{})

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(location) LIKE ? OR LOWER(description) LIKE ?",
			search, search, search,
		)
	}

	if filter.Location != "" {
		query = query.Where("LOWER(location) LIKE ?", "%"+strings.ToLower(filter.Location)+"%")
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := r.resolveZoneSortColumn(filter.SortBy)
	sortDir := "ASC"
	if strings.ToUpper(filter.SortDir) == "DESC" {
		sortDir = "DESC"
	}

	pagination := dto.NormalizePagination(filter.Page, filter.Limit)
	var zones []models.ParkingZone
	err := query.
		Order(sortBy + " " + sortDir).
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		Find(&zones).Error

	return zones, total, err
}

func (r *ParkingZoneRepository) resolveZoneSortColumn(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "name":
		return "name"
	case "location":
		return "location"
	case "capacity":
		return "capacity"
	case "hourly_rate":
		return "hourly_rate"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
