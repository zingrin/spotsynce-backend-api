package middleware

import (
	"net/http"
	"strings"

	apperrors "spotsync/pkg/errors"
	"spotsync/models"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

const (
	ContextUserIDKey = "user_id"
	ContextUserRole  = "user_role"
	ContextUserEmail = "user_email"
)

// JWTMiddleware validates JWT tokens and injects user claims into the context.
func JWTMiddleware(jwtManager *utils.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return apperrors.ErrUnauthorized
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return apperrors.ErrUnauthorized
			}

			claims, err := jwtManager.ValidateToken(parts[1])
			if err != nil {
				return apperrors.ErrUnauthorized
			}

			c.Set(ContextUserIDKey, claims.UserID)
			c.Set(ContextUserRole, claims.Role)
			c.Set(ContextUserEmail, claims.Email)

			return next(c)
		}
	}
}

// AdminMiddleware restricts access to admin users only.
func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get(ContextUserRole).(string)
			if !ok || role != models.RoleAdmin {
				return apperrors.ErrForbidden
			}
			return next(c)
		}
	}
}

// GetUserID extracts the authenticated user ID from the Echo context.
func GetUserID(c echo.Context) uint {
	userID, ok := c.Get(ContextUserIDKey).(uint)
	if !ok {
		return 0
	}
	return userID
}

// GetUserRole extracts the authenticated user role from the Echo context.
func GetUserRole(c echo.Context) string {
	role, ok := c.Get(ContextUserRole).(string)
	if !ok {
		return ""
	}
	return role
}

// CORS returns a middleware that handles Cross-Origin Resource Sharing headers.
func CORS() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Response().Header().Set("Access-Control-Max-Age", "86400")

			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}
