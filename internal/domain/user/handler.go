package user

import (
	"errors"
	"fmt"
	"net/http"
	"spot-sync/internal/domain/user/dto"
	"spot-sync/internal/httpresponse"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service}
}

func (h *handler) RegisterUser(c *echo.Context) error {
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

	res, err := h.service.RegisterUser(req)

	if err != nil {
		fmt.Println(err)

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, httpresponse.Response{
				Success: false,
				Message: "User already exists",
			})
		}

		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "User registration failed",
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    res,
	})
}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest

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

	res, err := h.service.LoginUser(req)

	if err != nil {
		fmt.Println(err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Response{
				Success: false,
				Message: "Invalid credentials",
			})
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Response{
				Success: false,
				Message: "Token expired",
			})
		}

		return c.JSON(http.StatusBadRequest, httpresponse.Response{
			Success: false,
			Message: "Login failed",
		})
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    res.Token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, httpresponse.Response{
		Success: true,
		Message: "Login successful",
		Data:    res,
	})
}
