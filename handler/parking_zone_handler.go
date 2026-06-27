package handler

import (
	"net/http"
	"strconv"

	"spotsync/dto"
	"spotsync/pkg/response"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

// ParkingZoneHandler handles parking zone HTTP requests.
type ParkingZoneHandler struct {
	zoneService *service.ParkingZoneService
}

// NewParkingZoneHandler creates a new ParkingZoneHandler.
func NewParkingZoneHandler(zoneService *service.ParkingZoneService) *ParkingZoneHandler {
	return &ParkingZoneHandler{zoneService: zoneService}
}

// Create godoc
// @Summary Create a parking zone
// @Description Create a new parking zone (admin only)
// @Tags Parking Zones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateZoneRequest true "Parking zone details"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/v1/zones [post]
func (h *ParkingZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if errors := utils.ValidateStruct(&req); errors != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", errors)
	}

	result, err := h.zoneService.Create(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Created(c, "parking zone created successfully", result)
}

// List godoc
// @Summary List parking zones
// @Description Get paginated list of parking zones with search, filtering, and sorting
// @Tags Parking Zones
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by name, location, or description"
// @Param location query string false "Filter by location"
// @Param is_active query bool false "Filter by active status"
// @Param sort_by query string false "Sort field (name, location, capacity, hourly_rate, created_at)" default(created_at)
// @Param sort_dir query string false "Sort direction (ASC, DESC)" default(DESC)
// @Success 200 {object} response.SuccessResponse
// @Router /api/v1/zones [get]
func (h *ParkingZoneHandler) List(c echo.Context) error {
	query := bindZoneListQuery(c)

	result, err := h.zoneService.List(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return response.OK(c, "parking zones retrieved successfully", result)
}

// GetByID godoc
// @Summary Get parking zone by ID
// @Description Retrieve a single parking zone by its ID
// @Tags Parking Zones
// @Produce json
// @Param id path int true "Parking zone ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/zones/{id} [get]
func (h *ParkingZoneHandler) GetByID(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid zone id")
	}

	result, err := h.zoneService.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, "parking zone retrieved successfully", result)
}

func bindZoneListQuery(c echo.Context) *dto.ZoneListQuery {
	query := &dto.ZoneListQuery{
		Page:     parseIntQuery(c, "page", 1),
		Limit:    parseIntQuery(c, "limit", 10),
		Search:   c.QueryParam("search"),
		Location: c.QueryParam("location"),
		SortBy:   c.QueryParam("sort_by"),
		SortDir:  c.QueryParam("sort_dir"),
	}

	if isActiveStr := c.QueryParam("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		query.IsActive = &isActive
	}

	return query
}

func parseIntQuery(c echo.Context, key string, defaultVal int) int {
	val := c.QueryParam(key)
	if val == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return parsed
}

func parseUintParam(c echo.Context, key string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
