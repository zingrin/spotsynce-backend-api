package middleware

import (
	"net/http"
	"spot-sync/internal/auth"
	"spot-sync/internal/httpresponse"
	"strings"

	"github.com/labstack/echo/v5"
)

const (
	ContextUserKey = "user"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, httpresponse.Response{
					Success: false,
					Message: "Unauthorized: Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, httpresponse.Response{
					Success: false,
					Message: "Unauthorized: Invalid authorization header",
				})
			}

			token := parts[1]

			claims, err := jwtService.ValidateToken(token)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, httpresponse.Response{
					Success: false,
					Message: "Invalid or expired token",
				})
			}

			c.Set(ContextUserKey, claims)

			return next(c)
		}
	}
}
