package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"goaconly/internal/infrastructure/config"
)

type Limiter interface {
	Allow(ctx context.Context, key string, limit int, windowSeconds int) (bool, int, error)
}

// RateLimitByIP используется для неаутентифицированных роутов (login, register),
// где userID из контекста ещё недоступен.
func RateLimitByIP(limiter Limiter, cfg config.RateLimitConfig, limit int) func(http.Handler) http.Handler {
	windowSeconds := int(cfg.WindowLength.Seconds())
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			key := fmt.Sprintf("ratelimit:ip:%s:%s", ip, r.URL.Path)

			allowed, current, err := limiter.Allow(r.Context(), key, limit, windowSeconds)
			if err != nil {
				// Fail-open: если Redis недоступен, не блокируем легитимный трафик,
				// но обязательно логируем — это сигнал деградации инфраструктуры.
				next.ServeHTTP(w, r)
				return
			}
			if !allowed {
				w.Header().Set("Retry-After", cfg.WindowLength.String())
				responseError(w, http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded: %d/%d", current, limit))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitByUser используется для защищённых роутов, применяется ПОСЛЕ
// middleware.Auth в цепочке — иначе userID в контексте ещё не будет.
func RateLimitByUser(limiter Limiter, cfg config.RateLimitConfig, limit int) func(http.Handler) http.Handler {
	windowSeconds := int(cfg.WindowLength.Seconds())
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := GetUserIDFromContext(r.Context())
			if err != nil {
				// Если это случилось — значит middleware подключён неправильно
				responseError(w, http.StatusInternalServerError, fmt.Errorf("rate limiter: user context missing"))
				return
			}
			key := fmt.Sprintf("ratelimit:user:%s:%s", userID.String(), r.URL.Path)

			allowed, current, err := limiter.Allow(r.Context(), key, limit, windowSeconds)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", cfg.WindowLength.String())
				responseError(w, http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded: %d/%d", current, limit))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// clientIP извлекает реальный IP клиента, учитывая возможный прокси/балансировщик
// перед приложением (в проде между интернетом и Go-сервисом почти всегда стоит nginx/ALB).
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.IndexByte(xff, ','); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
