package user

import "gorm.io/gorm"

type Repository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetByEmail(email string) (*User, error) {
	user := &User{}

	tx := r.db.Where(&User{Email: email}).First(user)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}
