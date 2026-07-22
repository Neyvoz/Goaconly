package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	usecaseAuth "goaconly/internal/usecase/auth"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "user_id"

func Auth(jwtService usecaseAuth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// достать заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responseError(w, http.StatusUnauthorized, errors.New("missing authorization header"))
				return
			}
			// проверить формат "Bearer <token>"
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				responseError(w, http.StatusUnauthorized, errors.New("invalid authorization header format"))
				return
			}
			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			// валидировать токен через jwtService.ValidateAccessToken
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				responseError(w, http.StatusUnauthorized, errors.New("invalid or expired token"))
				return
			}
			// положить claims.UserID в context
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user id not found in context")
	}
	return userID, nil
}
