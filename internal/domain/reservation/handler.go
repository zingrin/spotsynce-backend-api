package reservation

import (
	"errors"
	"fmt"
	"net/http"
	"spot-sync/internal/auth"
	"spot-sync/internal/domain/reservation/dto"
	"spot-sync/internal/domain/user"
	"spot-sync/internal/httpresponse"
	"spot-sync/internal/middleware"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service}
}

func (h *handler) ReserveSpot(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		fmt.Println(err)

		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		fmt.Println(err.Error())

		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Validation failed!",
			Error:   err.Error(),
		})
	}

	claims := c.Get(middleware.ContextUserKey).(*auth.JWTClaims)

	res, err := h.service.ReserveSpot(&req, claims.Id)

	if err != nil {
		switch {
		case errors.Is(err, ErrZoneFull):
			return c.JSON(http.StatusConflict, httpresponse.Response{
				Success: false,
				Message: "This zone is fully booked",
			})

		case errors.Is(err, ErrAlreadyReserved):
			return c.JSON(http.StatusConflict, httpresponse.Response{
				Success: false,
				Message: "This vehicle already has an active reservation",
			})

		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Response{
				Success: false,
				Message: "Failed to reserve spot",
			})
		}
	}

	return c.JSON(http.StatusCreated, httpresponse.Response{
		Success: true,
		Message: "Spot reserved successfully",
		Data:    res,
	})
}

func (h *handler) GetMyReservations(c *echo.Context) error {
	claims, ok := c.Get(middleware.ContextUserKey).(*auth.JWTClaims)

	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Response{
			Success: false,
			Message: "Unauthorized",
		})
	}

	res, err := h.service.GetMyReservations(claims.Id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Response{
			Success: false,
			Message: "Failed to retrieve reservations",
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    res,
	})
}

func (h *handler) CancelReservation(c *echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Invalid reservation id",
		})
	}

	claims, ok := c.Get(middleware.ContextUserKey).(*auth.JWTClaims)

	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Response{
			Success: false,
			Message: "Unauthorized",
		})
	}

	err = h.service.CancelReservation(id, claims.Id, user.UserRole(claims.Role))

	if err != nil {
		switch {
		case errors.Is(err, ErrReservationNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Response{
				Success: false,
				Message: "Reservation not found",
			})

		case errors.Is(err, ErrNotOwner):
			return c.JSON(http.StatusForbidden, httpresponse.Response{
				Success: false,
				Message: "You are not allowed to cancel this reservation",
			})

		case errors.Is(err, ErrInvalidStatusTransition):
			return c.JSON(http.StatusConflict, httpresponse.Response{
				Success: false,
				Message: "Only active reservations can be cancelled",
			})

		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Response{
				Success: false,
				Message: "Failed to cancel reservation",
			})
		}
	}

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

func (h *handler) GetAllReservations(c *echo.Context) error {
	res, err := h.service.GetAllReservations()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Response{
			Success: false,
			Message: "Failed to retrieve reservations",
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    res,
	})
}
