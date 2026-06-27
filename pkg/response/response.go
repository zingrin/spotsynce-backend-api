package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// SuccessResponse represents a standardized success JSON response.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standardized error JSON response.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// OK sends a 200 success response.
func OK(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 success response.
func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a 204 response with a success wrapper.
func NoContent(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
	})
}

// Error sends an error response with the given status code.
func Error(c echo.Context, statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
