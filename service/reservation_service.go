package service

import (
	"context"
	"errors"
	"math"
	"time"

	"spotsync/dto"
	apperrors "spotsync/pkg/errors"
	"spotsync/models"
	"spotsync/repository"

	"gorm.io/gorm"
)

// ReservationService handles reservation business logic.
type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	zoneRepo        *repository.ParkingZoneRepository
}

// NewReservationService creates a new ReservationService.
func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	zoneRepo *repository.ParkingZoneRepository,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

// Create books a parking spot with transactional row locking to prevent overbooking.
func (s *ReservationService) Create(ctx context.Context, userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	if req.StartTime.Before(time.Now()) {
		return nil, apperrors.NewWithDetails(400, "start time must be in the future", nil)
	}

	zone, err := s.zoneRepo.FindByID(ctx, req.ParkingZoneID)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to verify parking zone", nil)
	}
	if zone == nil {
		return nil, apperrors.ErrNotFound
	}
	if !zone.IsActive {
		return nil, apperrors.ErrZoneInactive
	}

	totalCost := calculateTotalCost(zone.HourlyRate, req.StartTime, req.EndTime)

	reservation := &models.Reservation{
		UserID:        userID,
		ParkingZoneID: req.ParkingZoneID,
		VehicleNumber: req.VehicleNumber,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		TotalCost:     totalCost,
		Status:        models.ReservationStatusActive,
	}

	if err := s.reservationRepo.CreateWithCapacityCheck(ctx, reservation); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, apperrors.ErrNotFound
		case errors.Is(err, repository.ErrZoneFull):
			return nil, apperrors.ErrZoneFull
		case errors.Is(err, repository.ErrZoneInactive):
			return nil, apperrors.ErrZoneInactive
		default:
			return nil, apperrors.NewWithDetails(500, "failed to create reservation", nil)
		}
	}

	resp := dto.ToReservationResponse(reservation)
	return &resp, nil
}

// GetMyReservations returns paginated reservations for the authenticated user.
func (s *ReservationService) GetMyReservations(ctx context.Context, userID uint, query *dto.ReservationListQuery) (*dto.PaginatedResponse, error) {
	filter := repository.ReservationFilter{
		UserID:  userID,
		Status:  query.Status,
		SortBy:  query.SortBy,
		SortDir: query.SortDir,
		Page:    query.Page,
		Limit:   query.Limit,
	}

	reservations, total, err := s.reservationRepo.List(ctx, filter, true)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to retrieve reservations", nil)
	}

	pagination := dto.NormalizePagination(query.Page, query.Limit)
	return &dto.PaginatedResponse{
		Items:      dto.ToReservationResponses(reservations),
		Total:      total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: dto.TotalPages(total, pagination.Limit),
	}, nil
}

// Cancel soft-deletes a reservation owned by the authenticated user.
func (s *ReservationService) Cancel(ctx context.Context, userID, reservationID uint) error {
	reservation, err := s.reservationRepo.FindByID(ctx, reservationID, false)
	if err != nil {
		return apperrors.NewWithDetails(500, "failed to retrieve reservation", nil)
	}
	if reservation == nil {
		return apperrors.ErrReservationNotFound
	}
	if reservation.UserID != userID {
		return apperrors.ErrForbidden
	}

	if err := s.reservationRepo.SoftDelete(ctx, reservationID); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return apperrors.ErrReservationNotFound
		case errors.Is(err, repository.ErrCannotCancel):
			return apperrors.ErrCannotCancel
		default:
			return apperrors.NewWithDetails(500, "failed to cancel reservation", nil)
		}
	}

	return nil
}

// ListAll returns all reservations with filtering (admin only).
func (s *ReservationService) ListAll(ctx context.Context, query *dto.ReservationListQuery) (*dto.PaginatedResponse, error) {
	filter := repository.ReservationFilter{
		UserID:        query.UserID,
		ParkingZoneID: query.ParkingZoneID,
		Status:        query.Status,
		SortBy:        query.SortBy,
		SortDir:       query.SortDir,
		Page:          query.Page,
		Limit:         query.Limit,
	}

	reservations, total, err := s.reservationRepo.List(ctx, filter, true)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to list reservations", nil)
	}

	pagination := dto.NormalizePagination(query.Page, query.Limit)
	return &dto.PaginatedResponse{
		Items:      dto.ToReservationResponses(reservations),
		Total:      total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: dto.TotalPages(total, pagination.Limit),
	}, nil
}

// calculateTotalCost computes reservation cost based on hourly rate and duration.
func calculateTotalCost(hourlyRate float64, start, end time.Time) float64 {
	duration := end.Sub(start).Hours()
	cost := hourlyRate * duration
	return math.Round(cost*100) / 100
}
