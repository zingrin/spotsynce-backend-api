package zone

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

	api.POST("/zones", handler.CreateZone, middleware.AuthMiddleware(jwt),
		middleware.RequireRole(user.ADMIN))
	api.GET("/zones", handler.GetAllZones)
	api.GET("/zones/:id", handler.GetZoneById)
	api.PATCH("/zones/:id", handler.UpdateZone, middleware.AuthMiddleware(jwt), middleware.RequireRole(user.ADMIN))
	api.DELETE("/zones/:id", handler.DeleteZone, middleware.AuthMiddleware(jwt), middleware.RequireRole(user.ADMIN))
}
