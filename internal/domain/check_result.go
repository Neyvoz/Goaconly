package domain

import "time"

// CheckResult - результат одной проверки
type CheckResult struct {
	TargetID       int64
	CheckedAt      time.Time
	StatusCode     int
	ResponseTimeMs int64
	IsUp           bool
	SSLExpiresAt   *time.Time
	SSLDaysLeft    *int
	KeywordFound   *bool
	ErrorMessage   string
}
