package server

import (
	"net/http"
	"spot-sync/internal/config"
	"spot-sync/internal/domain/reservation"
	"spot-sync/internal/domain/user"
	"spot-sync/internal/domain/zone"
	"spot-sync/internal/httpresponse"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB, env *config.Env) {
	// Create enums
	(func() {
		db.Exec(`
    		DO $$
    		BEGIN
        		IF NOT EXISTS (
            		SELECT 1 FROM pg_type WHERE typname = 'user_role'
        		) THEN
            		CREATE TYPE user_role AS ENUM (
                		'admin',
                		'driver'
            		);
        		ELSE
            		ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'admin';
            		ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'driver';
        		END IF;
    		END
    		$$;
		`)

		db.Exec(`
    		DO $$
    		BEGIN
        		IF NOT EXISTS (
            		SELECT 1 FROM pg_type WHERE typname = 'zone_type'
        		) THEN
            		CREATE TYPE zone_type AS ENUM (
                		'general',
                		'ev_charging',
                		'covered'
            		);
        		ELSE
            		ALTER TYPE zone_type ADD VALUE IF NOT EXISTS 'general';
            		ALTER TYPE zone_type ADD VALUE IF NOT EXISTS 'ev_charging';
            		ALTER TYPE zone_type ADD VALUE IF NOT EXISTS 'covered';
        		END IF;
    		END
    		$$;
		`)

		db.Exec(`
    		DO $$
    		BEGIN
        		IF NOT EXISTS (
            		SELECT 1 FROM pg_type WHERE typname = 'reservation_status'
        		) THEN
            		CREATE TYPE reservation_status AS ENUM (
                		'ACTIVE',
                		'COMPLETED',
                		'CANCELED'
            		);
        		ELSE
            		ALTER TYPE reservation_status ADD VALUE IF NOT EXISTS 'ACTIVE';
            		ALTER TYPE reservation_status ADD VALUE IF NOT EXISTS 'COMPLETED';
            		ALTER TYPE reservation_status ADD VALUE IF NOT EXISTS 'CANCELED';
        		END IF;
    		END
    		$$;
		`)
	})()

	db.AutoMigrate(
		&user.User{},
		&zone.Zone{},
		&reservation.Reservation{},
	)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, httpresponse.Response{
			Success: true,
			Message: "Server is running...",
		})
	})

	apiV1 := e.Group("/api/v1")

	user.Routes(db, apiV1, env)
	zone.Routes(db, apiV1, env)
	reservation.Routes(db, apiV1, env)
}
