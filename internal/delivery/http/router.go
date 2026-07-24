package httpserver

import (
	"goaconly/internal/delivery/http/handler"
	"goaconly/internal/delivery/http/middleware"
	"goaconly/internal/infrastructure/config"
	"net/http"

	usecaseAuth "goaconly/internal/usecase/auth"
)

type Dependencies struct {
	TargetHandler *handler.TargetHandler
	AuthHandler   *handler.AuthHandler
	JWTService    usecaseAuth.JWTService
	Limiter       middleware.Limiter
	RateLimitCfg  config.RateLimitConfig
}

func NewRouter(deps Dependencies) *http.ServeMux {
	mux := http.NewServeMux()
	protected := middleware.Auth(deps.JWTService)

	authLimiter := middleware.RateLimitByIP(deps.Limiter, deps.RateLimitCfg, deps.RateLimitCfg.AuthRPS)
	apiLimiter := middleware.RateLimitByUser(deps.Limiter, deps.RateLimitCfg, deps.RateLimitCfg.DefaultRPS)

	// Публичные маршруты — только IP-based rate limit
	mux.Handle("POST /api/v1/auth/register", authLimiter(http.HandlerFunc(deps.AuthHandler.Register)))
	mux.Handle("POST /api/v1/auth/login", authLimiter(http.HandlerFunc(deps.AuthHandler.Login)))
	mux.Handle("POST /api/v1/auth/refresh", authLimiter(http.HandlerFunc(deps.AuthHandler.Refresh)))
	mux.Handle("POST /api/v1/auth/logout", authLimiter(http.HandlerFunc(deps.AuthHandler.Logout)))
	// Защищённые маршруты — сначала Auth (кладёт userID в контекст), потом rate-limit по userID
	mux.Handle("GET /api/v1/targets", protected(apiLimiter(http.HandlerFunc(deps.TargetHandler.List))))
	mux.Handle("POST /api/v1/targets", protected(apiLimiter(http.HandlerFunc(deps.TargetHandler.Create))))
	mux.Handle("GET /api/v1/targets/{id}", protected(apiLimiter(http.HandlerFunc(deps.TargetHandler.GetByID))))
	mux.Handle("PUT /api/v1/targets/{id}", protected(apiLimiter(http.HandlerFunc(deps.TargetHandler.Update))))
	mux.Handle("DELETE /api/v1/targets/{id}", protected(apiLimiter(http.HandlerFunc(deps.TargetHandler.Delete))))

	return mux
}
