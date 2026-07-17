package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims — данные, зашитые в access-токен.
type Claims struct {
	jwt.RegisteredClaims
	UserID   uuid.UUID
	PlanType string
}

type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, planType string) (string, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
}
