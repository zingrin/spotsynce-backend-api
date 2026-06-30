package user

import (
	"spot-sync/internal/domain/user/dto"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	ADMIN  UserRole = "admin"
	DRIVER UserRole = "driver"
)

type User struct {
	Id        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	Email     string         `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password  string         `json:"-" gorm:"type:text;not null"`
	Role      UserRole       `json:"role" gorm:"type:user_role;not null"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"type:timestamp;index"`
}

func (u *User) toResponse() *dto.UserResponse {
	return &dto.UserResponse{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) hashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

func (u *User) checkPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		return err
	}

	return nil
}
