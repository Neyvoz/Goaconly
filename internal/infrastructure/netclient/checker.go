package netclient

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"goaconly/internal/domain"
)

const defaultScanLimit = 1 << 20 // 1MB

type Checker struct {
	client *http.Client
	logger *slog.Logger
}

func NewChecker(logger *slog.Logger) *Checker {
	return &Checker{
		client: NewHTTPClient(DefaultClientConfig()),
		logger: logger,
	}
}

func (c *Checker) Check(ctx context.Context, target domain.Target) (domain.CheckResult, error) {
	result := domain.CheckResult{
		TargetID:  target.ID,
		CheckedAt: time.Now(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.URL, nil)
	if err != nil {
		return result, fmt.Errorf("build request: %w", err)
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	result.ResponseTimeMs = time.Since(start).Milliseconds()

	if err != nil {
		result.IsUp = false
		result.ErrorMessage = err.Error()
		return result, nil // сетевая ошибка — валидный результат
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

	if sslInfo, err := ExtractSSLInfo(resp); err != nil {
		c.logger.Warn("ssl extraction failed", "url", target.URL, "error", err)
	} else if sslInfo != nil {
		result.SSLExpiresAt = &sslInfo.ExpiresAt
		result.SSLDaysLeft = &sslInfo.DaysLeft
	}

	if target.KeywordToFind != "" {
		found, err := ScanKeyword(resp.Body, target.KeywordToFind, defaultScanLimit)
		if err != nil {
			c.logger.Warn("keyword scan failed", "url", target.URL, "error", err)
		}
		result.KeywordFound = &found
	}

	return result, nil
}
