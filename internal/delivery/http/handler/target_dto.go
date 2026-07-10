package handler

import (
	"sitepulse/internal/domain"
	"time"
)

// CreateTargetRequest — DTO входящего запроса на создание цели мониторинга.
// Специально отделён от domain.Target, чтобы клиент API не был завязан
// на внутреннее устройство доменной модели.
type CreateTargetRequest struct {
	URL              string `json:"url"`
	KeywordToFind    string `json:"keyword"`
	CheckIntervalMin int    `json:"interval_seconds"`
}

// TargetResponse — DTO исходящего ответа. Interval отдаём в секундах (int),
// а не как time.Duration, потому что Duration сериализуется в JSON
// как число наносекунд — нечитаемо и неудобно для клиентов API.
type TargetResponse struct {
	ID               int64  `json:"id"`
	URL              string `json:"url"`
	KeywordToFind    string `json:"keyword"`
	CheckIntervalMin int    `json:"interval_seconds"`
	IsActive         bool   `json:"is_active"`
	CreatedAt        string `json:"created_at"`
}

// toTargetResponse маппит доменную модель в DTO ответа.
// Это единственное место, где домен "встречается" с представлением API —
// вся остальная система ничего не знает про JSON и HTTP.
func toTargetResponse(t domain.Target) TargetResponse {
	return TargetResponse{
		ID:               t.ID,
		URL:              t.URL,
		KeywordToFind:    t.KeywordToFind,
		CheckIntervalMin: int(t.CheckInterval),
		IsActive:         t.IsActive,
		CreatedAt:        t.CreatedAt.Format(time.RFC3339),
	}
}
