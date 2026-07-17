package domain

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken — доменное представление refresh-токена.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpireAt  time.Time
	Revoked   bool
	CreatedAt time.Time
}
