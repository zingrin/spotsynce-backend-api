// @title SpotSync API
// @version 1.0
// @description SpotSync is a production-ready parking reservation REST API built with Go, Echo, and PostgreSQL.
// @termsOfService http://swagger.io/terms/

// @contact.name SpotSync Support
// @contact.email support@spotsync.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"spotsync/config"
	"spotsync/database"
	_ "spotsync/docs"
	"spotsync/handler"
	"spotsync/repository"
	"spotsync/routes"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Dependency injection: Repository → Service → Handler
	jwtManager := utils.NewJWTManager(cfg)

	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewParkingZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo, jwtManager)
	zoneService := service.NewParkingZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo, zoneRepo)

	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewParkingZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)
	healthHandler := handler.NewHealthHandler(db)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = handler.ErrorHandler

	routes.Register(e, &routes.Handlers{
		Auth:        authHandler,
		Zone:        zoneHandler,
		Reservation: reservationHandler,
		Health:      healthHandler,
	}, jwtManager)

	// Graceful shutdown
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("SpotSync server starting on %s", addr)
		log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", cfg.Port)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}

	log.Println("server stopped")
}
