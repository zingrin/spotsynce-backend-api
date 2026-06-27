package handler

import (
	"net/http"

	"spotsync/dto"
	"spotsync/pkg/response"
	"spotsync/service"
	"spotsync/utils"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account and receive a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if errors := utils.ValidateStruct(&req); errors != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", errors)
	}

	result, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Created(c, "registration successful", result)
}

// Login godoc
// @Summary Login user
// @Description Authenticate with email and password to receive a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if errors := utils.ValidateStruct(&req); errors != nil {
		return response.Error(c, http.StatusBadRequest, "validation failed", errors)
	}

	result, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.OK(c, "login successful", result)
}
