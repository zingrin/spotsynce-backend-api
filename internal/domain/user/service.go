package user

import (
	"spot-sync/internal/auth"
	"spot-sync/internal/domain/user/dto"
)

type service struct {
	repo Repository
	jwt  auth.JWTService
}

func NewService(repo Repository, jwt auth.JWTService) *service {
	return &service{repo, jwt}
}

func (s *service) RegisterUser(req dto.CreateRequest) (*dto.UserResponse, error) {
	user := &User{
		Name:  req.Name,
		Email: req.Email,
		Role:  UserRole(req.Role),
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user.toResponse(), nil
}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	if err := user.checkPassword(req.Password); err != nil {
		return nil, err
	}

	token, err := s.jwt.GenerateToken(user.Id, user.Name, user.Email, string(user.Role))

	if err != nil {
		return nil, err
	}

	res := &dto.LoginResponse{
		Token: token,
		User:  *user.toResponse(),
	}

	return res, nil
}
