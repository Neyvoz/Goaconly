package usecase

import (
	"context"
	"sitepulse/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email domain.Email) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type RefreshTokenRepository interface {
	Store(ctx context.Context, t domain.RefreshToken) error
	GetByHash(ctx context.Context) (domain.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
