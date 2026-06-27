package utils

import (
	"fmt"
	"time"

	"spotsync/config"
	"spotsync/models"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims holds the custom JWT claims for authenticated users.
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation.
type JWTManager struct {
	secret     []byte
	expiryTime time.Duration
}

// NewJWTManager creates a new JWTManager from application config.
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{
		secret:     []byte(cfg.JWTSecret),
		expiryTime: cfg.JWTExpiryDuration(),
	}
}

// GenerateToken creates a signed JWT for the given user.
func (m *JWTManager) GenerateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiryTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken parses and validates a JWT string, returning the claims.
func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
