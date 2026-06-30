package zone

import (
	"errors"
	"fmt"
	"net/http"
	"spot-sync/internal/domain/zone/dto"
	"spot-sync/internal/httpresponse"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service}
}

func (h *handler) CreateZone(c *echo.Context) error {
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

	res, err := h.service.CreateZone(req)

	if err != nil {
		fmt.Println(err)

		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Failed to create parking zone",
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Response{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    res,
	})
}

func (h *handler) GetAllZones(c *echo.Context) error {
	res, err := h.service.GetAllZones()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Response{
			Success: false,
			Message: "Failed to retrieve parking zones",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    res,
	})
}

func (h *handler) GetZoneById(c *echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Invalid zone id",
		})
	}

	res, err := h.service.GetZoneById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Response{
				Success: false,
				Message: "Zone not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Response{
			Success: false,
			Message: "Failed to get zone",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "Zone retrieved successfully",
		Data:    res,
	})
}

func (h *handler) UpdateZone(c *echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Invalid zone id",
		})
	}

	var req dto.UpdateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Validation failed!",
			Error:   err.Error(),
		})
	}

	res, err := h.service.UpdateZone(id, &req)

	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Response{
				Success: false,
				Message: "Zone not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Response{
			Success: false,
			Message: "Failed to update zone",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "Zone updated successfully",
		Data:    res,
	})
}

func (h *handler) DeleteZone(c *echo.Context) error {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Invalid zone id",
		})
	}

	if err := h.service.DeleteZone(id); err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Response{
				Success: false,
				Message: "Zone not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Response{
			Success: false,
			Message: "Failed to delete zone",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "Zone deleted successfully",
	})
}
