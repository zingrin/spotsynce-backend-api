package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"spotsync/dto"
	"spotsync/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReservationFilter holds filter criteria for reservation queries.
type ReservationFilter struct {
	UserID        uint
	ParkingZoneID uint
	Status        string
	SortBy        string
	SortDir       string
	Page          int
	Limit         int
}

// ReservationRepository handles all database operations for reservations.
type ReservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new ReservationRepository.
func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// CreateWithCapacityCheck atomically locks the parking zone, counts overlapping
// active reservations, and creates a new reservation if capacity allows.
func (r *ReservationRepository) CreateWithCapacityCheck(
	ctx context.Context,
	reservation *models.Reservation,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, reservation.ParkingZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			return err
		}

		if !zone.IsActive {
			return ErrZoneInactive
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("parking_zone_id = ?", zone.ID).
			Where("status = ?", models.ReservationStatusActive).
			Where("start_time < ? AND end_time > ?", reservation.EndTime, reservation.StartTime).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.Capacity {
			return ErrZoneFull
		}

		return tx.Create(reservation).Error
	})
}

// FindByID retrieves a reservation by primary key with optional preloads.
func (r *ReservationRepository) FindByID(ctx context.Context, id uint, preload bool) (*models.Reservation, error) {
	query := r.db.WithContext(ctx)
	if preload {
		query = query.Preload("ParkingZone").Preload("User")
	}

	var reservation models.Reservation
	err := query.First(&reservation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reservation, nil
}

// List retrieves reservations with filtering, sorting, and pagination.
func (r *ReservationRepository) List(ctx context.Context, filter ReservationFilter, preload bool) ([]models.Reservation, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Reservation{})

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.ParkingZoneID > 0 {
		query = query.Where("parking_zone_id = ?", filter.ParkingZoneID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := r.resolveReservationSortColumn(filter.SortBy)
	sortDir := "DESC"
	if strings.ToUpper(filter.SortDir) == "ASC" {
		sortDir = "ASC"
	}

	pagination := dto.NormalizePagination(filter.Page, filter.Limit)

	findQuery := r.db.WithContext(ctx).Model(&models.Reservation{})
	if filter.UserID > 0 {
		findQuery = findQuery.Where("user_id = ?", filter.UserID)
	}
	if filter.ParkingZoneID > 0 {
		findQuery = findQuery.Where("parking_zone_id = ?", filter.ParkingZoneID)
	}
	if filter.Status != "" {
		findQuery = findQuery.Where("status = ?", filter.Status)
	}
	if preload {
		findQuery = findQuery.Preload("ParkingZone").Preload("User")
	}

	var reservations []models.Reservation
	err := findQuery.
		Order(sortBy + " " + sortDir).
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		Find(&reservations).Error

	return reservations, total, err
}

// SoftDelete cancels a reservation by setting status and soft-deleting the record.
func (r *ReservationRepository) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reservation models.Reservation
		if err := tx.First(&reservation, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			return err
		}

		if reservation.Status != models.ReservationStatusActive {
			return ErrCannotCancel
		}

		if reservation.EndTime.Before(time.Now()) {
			return ErrCannotCancel
		}

		reservation.Status = models.ReservationStatusCancelled
		if err := tx.Save(&reservation).Error; err != nil {
			return err
		}

		return tx.Delete(&reservation).Error
	})
}

func (r *ReservationRepository) resolveReservationSortColumn(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "start_time":
		return "start_time"
	case "end_time":
		return "end_time"
	case "total_cost":
		return "total_cost"
	case "status":
		return "status"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}

// Repository-level sentinel errors mapped by the service layer.
var (
	ErrZoneFull      = errors.New("zone full")
	ErrZoneInactive  = errors.New("zone inactive")
	ErrCannotCancel  = errors.New("cannot cancel")
)
