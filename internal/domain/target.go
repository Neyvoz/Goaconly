package domain

import (
	"time"

	"github.com/google/uuid"
)

// CheckInterval — типизированный алиас для предотвращения путаницы с int
type CheckInterval int

// CheckJob — единица работы, которую воркер получает из канала.
type CheckJob struct {
	Target Target
}

// Target — сайт или endpoint, который система мониторит.
// Это центральная сущность домена.
type Target struct {
	ID            int64
	UserID        uuid.UUID
	URL           string
	CheckInterval CheckInterval // в минутах
	KeywordToFind string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
