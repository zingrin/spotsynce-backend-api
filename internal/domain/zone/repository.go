package zone

import (
	"spot-sync/internal/domain/zone/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(zone *Zone) error
	FindAllWithAvailability() ([]dto.ZoneWithAvailability, error)
	FindByIdWithAvailability(id uuid.UUID) (*dto.ZoneWithAvailability, error)
	FindById(id uuid.UUID) (*Zone, error)
	Update(zone *Zone) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(zone *Zone) error {
	return r.db.Create(zone).Error
}

func (r *repository) FindAllWithAvailability() ([]dto.ZoneWithAvailability, error) {
	var zones []dto.ZoneWithAvailability

	err := r.db.Table("zones AS z").
		Select(`
			z.id,
			z.name,
			z.type,
			z.total_capacity,
			z.total_capacity - COALESCE((
				SELECT COUNT(*) FROM reservations r
				WHERE r.zone_id = z.id
				AND r.status = ?
				AND r.deleted_at IS NULL
			), 0) AS available_spots,
			z.price_per_hour,
			z.created_at
		`, "ACTIVE").
		Where("z.deleted_at IS NULL").
		Find(&zones).Error

	if err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *repository) FindByIdWithAvailability(id uuid.UUID) (*dto.ZoneWithAvailability, error) {
	var zone dto.ZoneWithAvailability

	err := r.db.Table("zones AS z").
		Select(`
			z.id,
			z.name,
			z.type,
			z.total_capacity,
			z.total_capacity - COALESCE((
				SELECT COUNT(*) FROM reservations r
				WHERE r.zone_id = z.id
				AND r.status = ?
				AND r.deleted_at IS NULL
			), 0) AS available_spots,
			z.price_per_hour,
			z.created_at
		`, "ACTIVE").
		Where("z.id = ? AND z.deleted_at IS NULL", id).
		First(&zone).Error

	if err != nil {
		return nil, err
	}

	return &zone, nil
}

func (r *repository) FindById(id uuid.UUID) (*Zone, error) {
	var zone Zone

	if err := r.db.First(&zone, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &zone, nil
}

func (r *repository) Update(zone *Zone) error {
	return r.db.
		Model(&Zone{}).
		Where(&Zone{Id: zone.Id}).
		Updates(zone).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.
		Model(&Zone{}).
		Where(&Zone{Id: id}).
		Delete(&Zone{}).Error
}
