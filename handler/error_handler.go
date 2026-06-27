package handler

import (
	"net/http"

	apperrors "spotsync/pkg/errors"
	"spotsync/pkg/response"

	"github.com/labstack/echo/v4"
)

// ErrorHandler is the centralized Echo HTTP error handler.
func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if appErr, ok := apperrors.IsAppError(err); ok {
		_ = response.Error(c, appErr.Code, appErr.Message, appErr.Errors)
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		message := "request error"
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
		_ = response.Error(c, he.Code, message, nil)
		return
	}

	_ = response.Error(c, http.StatusInternalServerError, "internal server error", nil)
}
