package handler

import (
	"net/http"
	"strconv"

	"spotsync/dto"
	"spotsync/middleware"
	"spotsync/pkg/response"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

// ReservationHandler handles reservation HTTP requests.
type ReservationHandler struct {
	reservationService *service.ReservationService
}

// NewReservationHandler creates a new ReservationHandler.
func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

// Create godoc
// @Summary Create a reservation
// @Description Book a parking spot in a zone (uses row locking to prevent overbooking)
// @Tags Reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateReservationRequest true "Reservation details"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/v1/reservations [post]
func (h *ReservationHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if errors := utils.ValidateStruct(&req); errors != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", errors)
	}

	result, err := h.reservationService.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return response.Created(c, "reservation created successfully", result)
}

// GetMyReservations godoc
// @Summary Get my reservations
// @Description Retrieve paginated reservations for the authenticated user
// @Tags Reservations
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status (active, cancelled, completed)"
// @Param sort_by query string false "Sort field (start_time, end_time, total_cost, status, created_at)"
// @Param sort_dir query string false "Sort direction (ASC, DESC)"
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/reservations/my-reservations [get]
func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID := middleware.GetUserID(c)
	query := bindMyReservationListQuery(c)

	result, err := h.reservationService.GetMyReservations(c.Request().Context(), userID, query)
	if err != nil {
		return err
	}

	return response.OK(c, "reservations retrieved successfully", result)
}

// Cancel godoc
// @Summary Cancel a reservation
// @Description Cancel and soft-delete a reservation owned by the authenticated user
// @Tags Reservations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/reservations/{id} [delete]
func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid reservation id")
	}

	if err := h.reservationService.Cancel(c.Request().Context(), userID, id); err != nil {
		return err
	}

	return response.NoContent(c, "reservation cancelled successfully")
}

// ListAll godoc
// @Summary List all reservations
// @Description Retrieve all reservations with filtering (admin only)
// @Tags Reservations
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status (active, cancelled, completed)"
// @Param parking_zone_id query int false "Filter by parking zone ID"
// @Param user_id query int false "Filter by user ID"
// @Param sort_by query string false "Sort field (start_time, end_time, total_cost, status, created_at)"
// @Param sort_dir query string false "Sort direction (ASC, DESC)"
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/v1/reservations [get]
func (h *ReservationHandler) ListAll(c echo.Context) error {
	query := bindReservationListQuery(c)

	result, err := h.reservationService.ListAll(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return response.OK(c, "reservations retrieved successfully", result)
}

func bindMyReservationListQuery(c echo.Context) *dto.ReservationListQuery {
	return &dto.ReservationListQuery{
		Page:    parseIntQuery(c, "page", 1),
		Limit:   parseIntQuery(c, "limit", 10),
		Status:  c.QueryParam("status"),
		SortBy:  c.QueryParam("sort_by"),
		SortDir: c.QueryParam("sort_dir"),
	}
}

func bindReservationListQuery(c echo.Context) *dto.ReservationListQuery {
	query := &dto.ReservationListQuery{
		Page:    parseIntQuery(c, "page", 1),
		Limit:   parseIntQuery(c, "limit", 10),
		Status:  c.QueryParam("status"),
		SortBy:  c.QueryParam("sort_by"),
		SortDir: c.QueryParam("sort_dir"),
	}

	if zoneIDStr := c.QueryParam("parking_zone_id"); zoneIDStr != "" {
		if id, err := parseUintParamFromString(zoneIDStr); err == nil {
			query.ParkingZoneID = id
		}
	}

	if userIDStr := c.QueryParam("user_id"); userIDStr != "" {
		if id, err := parseUintParamFromString(userIDStr); err == nil {
			query.UserID = id
		}
	}

	return query
}

func parseUintParamFromString(val string) (uint, error) {
	parsed, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
