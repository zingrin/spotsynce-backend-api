package middleware

import (
	"net/http"
	"slices"
	"spot-sync/internal/auth"
	"spot-sync/internal/domain/user"
	"spot-sync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

func RequireRole(roles ...user.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			claimsRaw := c.Get(ContextUserKey)

			claims, ok := claimsRaw.(*auth.JWTClaims)

			if !ok {
				return c.JSON(http.StatusUnauthorized, httpresponse.Response{
					Success: false,
					Message: "Unauthorized",
				})
			}

			if slices.Contains(roles, user.UserRole(claims.Role)) {
				return next(c)
			}

			return c.JSON(http.StatusForbidden, httpresponse.Response{
				Success: false,
				Message: "You do not have permission to access this resource",
			})
		}
	}
}
