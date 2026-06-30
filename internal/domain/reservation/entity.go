package reservation

import (
	"spot-sync/internal/domain/reservation/dto"
	"spot-sync/internal/domain/user"
	"spot-sync/internal/domain/zone"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationStatus string

const (
	ACTIVE    ReservationStatus = "ACTIVE"
	COMPLETED ReservationStatus = "COMPLETED"
	CANCELED  ReservationStatus = "CANCELED"
)

type Reservation struct {
	Id           uuid.UUID         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	UserId       uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	User         user.User         `json:"user,omitempty" gorm:"foreignKey:UserId"`
	ZoneId       uuid.UUID         `json:"zone_id" gorm:"type:uuid;not null"`
	Zone         zone.Zone         `json:"zone,omitempty" gorm:"foreignKey:ZoneId"`
	LicensePlate string            `json:"license_plate" gorm:"type:varchar(15);index;not null"`
	Status       ReservationStatus `json:"status" gorm:"type:reservation_status;default:'ACTIVE';not null"`
	CreatedAt    time.Time         `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt    gorm.DeletedAt    `json:"-" gorm:"type:timestamp;index"`
}

func (r *Reservation) toResponse() *dto.ReservationResponse {
	return &dto.ReservationResponse{
		Id:           r.Id,
		UserId:       r.UserId,
		ZoneId:       r.ZoneId,
		LicensePlate: r.LicensePlate,
		Status:       string(r.Status),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func (r *Reservation) toMyResponse() *dto.MyReservationResponse {
	return &dto.MyReservationResponse{
		Id:           r.Id,
		LicensePlate: r.LicensePlate,
		Status:       string(r.Status),
		Zone: dto.ZoneInfo{
			Id:   r.Zone.Id,
			Name: r.Zone.Name,
			Type: string(r.Zone.Type),
		},
		CreatedAt: r.CreatedAt,
	}
}

func (r *Reservation) toAdminResponse() *dto.AdminReservationResponse {
	return &dto.AdminReservationResponse{
		Id:           r.Id,
		LicensePlate: r.LicensePlate,
		Status:       string(r.Status),
		User: dto.UserInfo{
			Id:    r.User.Id,
			Name:  r.User.Name,
			Email: r.User.Email,
		},
		Zone: dto.ZoneInfo{
			Id:   r.Zone.Id,
			Name: r.Zone.Name,
			Type: string(r.Zone.Type),
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
