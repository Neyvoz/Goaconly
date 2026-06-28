package domain

import "time"

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
	UserID        int64
	URL           string
	CheckInterval CheckInterval // в минутах
	KeywordToFind string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
