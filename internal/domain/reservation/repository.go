package reservation

import (
	"errors"
	"spot-sync/internal/domain/zone"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(reservation *Reservation) error
	FindByUserId(userId uuid.UUID) ([]Reservation, error)
	FindById(id uuid.UUID) (*Reservation, error)
	UpdateStatus(id uuid.UUID, status ReservationStatus) error
	GetAll() ([]Reservation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

var (
	ErrZoneFull                = errors.New("zone is fully booked")
	ErrAlreadyReserved         = errors.New("license plate already has an active reservation")
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrNotOwner                = errors.New("you are not allowed to cancel this reservation")
	ErrInvalidStatusTransition = errors.New("only active reservations can be cancelled")
)

func (r *repository) Create(reservation *Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var z zone.Zone

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(&zone.Zone{Id: reservation.Id}).
			First(&z).
			Error; err != nil {
			return err
		}

		var activeCount int64

		if err := tx.
			Model(&Reservation{}).
			Where(&Reservation{ZoneId: z.Id, Status: ACTIVE}).
			Count(&activeCount).
			Error; err != nil {
			return err
		}

		if activeCount >= int64(z.TotalCapacity) {
			return ErrZoneFull
		}

		var existing int64

		if err := tx.
			Model(&Reservation{}).
			Where(&Reservation{LicensePlate: reservation.LicensePlate, Status: ACTIVE}).
			Count(&existing).
			Error; err != nil {
			return err
		}

		if existing > 0 {
			return ErrAlreadyReserved
		}

		return tx.Create(reservation).Error
	})
}

func (r *repository) FindByUserId(userId uuid.UUID) ([]Reservation, error) {
	var reservations []Reservation

	if err := r.db.
		Preload("Zone").
		Where(&Reservation{UserId: userId}).
		Order("created_at desc").
		Find(&reservations).
		Error; err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *repository) FindById(id uuid.UUID) (*Reservation, error) {
	var reservation Reservation

	if err := r.db.
		Where(&Reservation{Id: id}).
		First(&reservation).
		Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *repository) UpdateStatus(id uuid.UUID, status ReservationStatus) error {
	return r.db.
		Model(&Reservation{}).
		Where(&Reservation{Id: id}).
		Update("status", status).Error
}

func (r *repository) GetAll() ([]Reservation, error) {
	var reservations []Reservation

	if err := r.db.
		Preload("User").
		Preload("Zone").
		Order("created_at desc").
		Find(&reservations).
		Error; err != nil {
		return nil, err
	}

	return reservations, nil
}
