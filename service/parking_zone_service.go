package service

import (
	"context"

	"spotsync/dto"
	apperrors "spotsync/pkg/errors"
	"spotsync/models"
	"spotsync/repository"
)

// ParkingZoneService handles parking zone business logic.
type ParkingZoneService struct {
	zoneRepo *repository.ParkingZoneRepository
}

// NewParkingZoneService creates a new ParkingZoneService.
func NewParkingZoneService(zoneRepo *repository.ParkingZoneRepository) *ParkingZoneService {
	return &ParkingZoneService{zoneRepo: zoneRepo}
}

// Create adds a new parking zone (admin only).
func (s *ParkingZoneService) Create(ctx context.Context, req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
		Capacity:    req.Capacity,
		HourlyRate:  req.HourlyRate,
		IsActive:    true,
	}

	if err := s.zoneRepo.Create(ctx, zone); err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to create parking zone", nil)
	}

	resp := dto.ToZoneResponse(zone)
	return &resp, nil
}

// GetByID retrieves a single parking zone by ID.
func (s *ParkingZoneService) GetByID(ctx context.Context, id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to retrieve parking zone", nil)
	}
	if zone == nil {
		return nil, apperrors.ErrNotFound
	}

	resp := dto.ToZoneResponse(zone)
	return &resp, nil
}

// List retrieves parking zones with pagination, filtering, sorting, and search.
func (s *ParkingZoneService) List(ctx context.Context, query *dto.ZoneListQuery) (*dto.PaginatedResponse, error) {
	filter := repository.ZoneFilter{
		Search:   query.Search,
		Location: query.Location,
		IsActive: query.IsActive,
		SortBy:   query.SortBy,
		SortDir:  query.SortDir,
		Page:     query.Page,
		Limit:    query.Limit,
	}

	zones, total, err := s.zoneRepo.List(ctx, filter)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to list parking zones", nil)
	}

	pagination := dto.NormalizePagination(query.Page, query.Limit)
	return &dto.PaginatedResponse{
		Items:      dto.ToZoneResponses(zones),
		Total:      total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: dto.TotalPages(total, pagination.Limit),
	}, nil
}
