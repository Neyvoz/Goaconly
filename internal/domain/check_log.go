package domain

import "time"

// CheckStatus — тип статуса, использование строк вместо iota
// обоснованно: значения хранятся в БД и читаются людьми
type CheckStatus string

const (
	StatusUp      CheckStatus = "up"
	StatusDown    CheckStatus = "down"
	StatusPending CheckStatus = "pending_down"
)

// CheckLog — иммутабельная запись факта проверки.
// После создания не изменяется — только читается.
type CheckLog struct {
	ID           int64
	TargetID     int64
	CheckedAt    time.Time
	StatusCode   int
	ResponseTime time.Duration
	Status       CheckStatus
	ErrorMessage string     // пусто если Status == StatusUp
	SSLExpiresAt *time.Time // nil если сайт без HTTPS
}
