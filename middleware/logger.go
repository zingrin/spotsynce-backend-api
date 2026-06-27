package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Logger returns a configured request logger middleware.
func Logger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human} ${remote_ip}\n",
	})
}

// Recover returns a panic recovery middleware.
func Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}

// RequestID adds a unique request ID to each request for tracing.
func RequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return time.Now().Format("20060102150405.000000")
		},
	})
}
