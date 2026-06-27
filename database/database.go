package database

import (
	"fmt"
	"log"
	"time"

	"spotsync/config"
	"spotsync/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection establishes a PostgreSQL connection with connection pool settings.
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Production connection pool configuration
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	return db, nil
}

// AutoMigrate runs GORM auto-migration for all models.
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}
	log.Println("database migration completed successfully")
	return nil
}
