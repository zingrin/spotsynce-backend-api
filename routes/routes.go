package routes

import (
	"spotsync/handler"
	"spotsync/middleware"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Handlers groups all HTTP handlers for route registration.
type Handlers struct {
	Auth        *handler.AuthHandler
	Zone        *handler.ParkingZoneHandler
	Reservation *handler.ReservationHandler
	Health      *handler.HealthHandler
}

// Register sets up all application routes and middleware.
func Register(e *echo.Echo, h *Handlers, jwtManager *utils.JWTManager) {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	e.GET("/health", h.Health.Check)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")

	// Authentication routes (public)
	auth := v1.Group("/auth")
	auth.POST("/register", h.Auth.Register)
	auth.POST("/login", h.Auth.Login)

	// Parking zone routes
	zones := v1.Group("/zones")
	zones.GET("", h.Zone.List)
	zones.GET("/:id", h.Zone.GetByID)
	zones.POST("", h.Zone.Create, middleware.JWTMiddleware(jwtManager), middleware.AdminMiddleware())

	// Reservation routes
	reservations := v1.Group("/reservations")
	reservations.GET("", h.Reservation.ListAll, middleware.JWTMiddleware(jwtManager), middleware.AdminMiddleware())
	reservations.POST("", h.Reservation.Create, middleware.JWTMiddleware(jwtManager))
	reservations.GET("/my-reservations", h.Reservation.GetMyReservations, middleware.JWTMiddleware(jwtManager))
	reservations.DELETE("/:id", h.Reservation.Cancel, middleware.JWTMiddleware(jwtManager))
}
