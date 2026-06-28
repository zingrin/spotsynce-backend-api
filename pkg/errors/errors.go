package apperrors

import (
	"errors"
	"net/http"
)

// AppError represents a domain-level error with an HTTP status code
type AppError struct {
	Code    int
	Message string
	Errors  interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError.
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// NewWithDetails creates a new AppError with additional error details.
func NewWithDetails(code int, message string, details interface{}) *AppError {
	return &AppError{Code: code, Message: message, Errors: details}
}

// Common application errors
var (
	ErrNotFound            = New(http.StatusNotFound, "resource not found")
	ErrUnauthorized        = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = New(http.StatusForbidden, "forbidden")
	ErrConflict            = New(http.StatusConflict, "resource conflict")
	ErrBadRequest          = New(http.StatusBadRequest, "bad request")
	ErrInternalServer      = New(http.StatusInternalServerError, "internal server error")
	ErrInvalidCredentials  = New(http.StatusUnauthorized, "invalid email or password")
	ErrEmailAlreadyExists  = New(http.StatusConflict, "email already registered")
	ErrZoneFull            = New(http.StatusConflict, "parking zone is at full capacity for the selected time slot")
	ErrZoneInactive        = New(http.StatusBadRequest, "parking zone is not active")
	ErrReservationNotFound = New(http.StatusNotFound, "reservation not found")
	ErrCannotCancel        = New(http.StatusBadRequest, "reservation cannot be cancelled")
)

// IsAppError checks if an error is an AppError.
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
