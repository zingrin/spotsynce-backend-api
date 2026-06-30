package zone

import (
	"errors"
	"spot-sync/internal/domain/zone/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

var ErrZoneNotFound = errors.New("zone not found")

func (s *service) CreateZone(req dto.CreateRequest) (*dto.ZoneResponse, error) {
	zone := &Zone{
		Name:          req.Name,
		Type:          ZoneType(req.Type),
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, err
	}

	return zone.toResponse(), nil
}

func (s *service) GetAllZones() ([]dto.ZoneWithAvailability, error) {
	zones, err := s.repo.FindAllWithAvailability()

	if err != nil {
		return nil, err
	}

	res := make([]dto.ZoneWithAvailability, 0, len(zones))

	for _, z := range zones {
		res = append(res, dto.ZoneWithAvailability{
			Id:             z.Id,
			Name:           z.Name,
			Type:           z.Type,
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: z.AvailableSpots,
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt,
		})
	}

	return res, nil
}

func (s *service) GetZoneById(id uuid.UUID) (*dto.ZoneWithAvailability, error) {
	zone, err := s.repo.FindByIdWithAvailability(id)

	if err != nil {
		return nil, err
	}

	res := &dto.ZoneWithAvailability{
		Id:             zone.Id,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.AvailableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}

	return res, nil
}

func (s *service) UpdateZone(id uuid.UUID, req *dto.UpdateRequest) (*dto.ZoneResponse, error) {
	existing, err := s.repo.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}

		return nil, err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}

	if req.Type != "" {
		existing.Type = ZoneType(req.Type)
	}

	if req.TotalCapacity > 0 {
		existing.TotalCapacity = req.TotalCapacity
	}

	if req.PricePerHour > 0 {
		existing.PricePerHour = req.PricePerHour
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing.toResponse(), nil
}

func (s *service) DeleteZone(id uuid.UUID) error {
	_, err := s.repo.FindById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrZoneNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}
