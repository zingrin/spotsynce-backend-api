package handler

import (
	"net/http"

	"spotsync/pkg/response"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse holds health check result data.
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// Check godoc
// @Summary Health check
// @Description Check API and database connectivity status
// @Tags Health
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /health [get]
func (h *HealthHandler) Check(c echo.Context) error {
	dbStatus := "connected"

	sqlDB, err := h.db.DB()
	if err != nil {
		dbStatus = "disconnected"
		return response.Error(c, http.StatusServiceUnavailable, "service unavailable", HealthResponse{
			Status:   "unhealthy",
			Database: dbStatus,
		})
	}

	if err := sqlDB.Ping(); err != nil {
		dbStatus = "disconnected"
		return response.Error(c, http.StatusServiceUnavailable, "service unavailable", HealthResponse{
			Status:   "unhealthy",
			Database: dbStatus,
		})
	}

	return response.OK(c, "service is healthy", HealthResponse{
		Status:   "healthy",
		Database: dbStatus,
	})
}
