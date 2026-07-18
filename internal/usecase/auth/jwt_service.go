package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims — данные, зашитые в access-токен.
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
	PlanID int64     `json:"plan_id"`
}

type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, planID int64) (string, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
}
