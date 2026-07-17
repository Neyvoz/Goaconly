package auth

import (
	"errors"
	"fmt"
	"sitepulse/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	usecaseAuth "sitepulse/internal/usecase/auth"
)

// JWTManager - это менеджер JSON веб-токенов
type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTManager возвращает новый JWT менеджер
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

// GenerateAccessToken генерирует и подписывает новый токен для пользователя
func (m *JWTManager) GenerateAccessToken(userID uuid.UUID, planType string) (string, error) {
	now := time.Now()
	claims := usecaseAuth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenDuration)),
		},
		UserID:   userID,
		PlanType: planType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ValidateAccessToken проверяет строку с токеном доступа и возвращает UserClaims, если токен действителен
func (m *JWTManager) ValidateAccessToken(tokenString string) (*usecaseAuth.Claims, error) {
	claims := &usecaseAuth.Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*usecaseAuth.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
