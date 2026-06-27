package service

import (
	"context"

	"spotsync/dto"
	apperrors "spotsync/pkg/errors"
	"spotsync/models"
	"spotsync/repository"
	"spotsync/utils"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user account and returns an auth token.
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to check email availability", nil)
	}
	if exists {
		return nil, apperrors.ErrEmailAlreadyExists
	}

	role := models.RoleDriver
	if req.Role != "" {
		role = req.Role
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to process password", nil)
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to create user account", nil)
	}

	token, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to generate token", nil)
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "authentication failed", nil)
	}
	if user == nil || !utils.CheckPassword(req.Password, user.Password) {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, apperrors.NewWithDetails(500, "failed to generate token", nil)
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}
