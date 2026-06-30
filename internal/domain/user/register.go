package user

import (
	"spot-sync/internal/auth"
	"spot-sync/internal/config"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, api *echo.Group, env *config.Env) {
	repo := NewRepository(db)
	jwt := auth.NewJWTService(env.JWT_SECRET)
	service := NewService(repo, jwt)
	handler := NewHandler(service)

	api.POST("/auth/register", handler.RegisterUser)
	api.POST("/auth/login", handler.LoginUser)
}
