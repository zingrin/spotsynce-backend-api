package reservation

import (
	"spot-sync/internal/auth"
	"spot-sync/internal/config"
	"spot-sync/internal/domain/user"
	"spot-sync/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, api *echo.Group, env *config.Env) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)
	jwt := auth.NewJWTService(env.JWT_SECRET)

	api.POST("/reservations", handler.ReserveSpot, middleware.AuthMiddleware(jwt))
	api.GET("/reservations/my-reservations", handler.GetMyReservations, middleware.AuthMiddleware(jwt))
	api.DELETE("/reservations/:id", handler.CancelReservation, middleware.AuthMiddleware(jwt))
	api.GET("/reservations", handler.GetAllReservations, middleware.AuthMiddleware(jwt),
		middleware.RequireRole(user.ADMIN))
}
