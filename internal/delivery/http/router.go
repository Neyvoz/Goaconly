package httpserver

import (
	"goaconly/internal/delivery/http/handler"
	"goaconly/internal/delivery/http/middleware"
	"net/http"

	usecaseAuth "goaconly/internal/usecase/auth"
)

type Dependencies struct {
	TargetHandler *handler.TargetHandler
	AuthHandler   *handler.AuthHandler
	JWTService    usecaseAuth.JWTService
}

func NewRouter(deps Dependencies) *http.ServeMux {
	mux := http.NewServeMux()
	protected := middleware.Auth(deps.JWTService)

	// Публичные маршруты — без middleware
	mux.HandleFunc("POST /api/v1/auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", deps.AuthHandler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", deps.AuthHandler.Logout)
	// Защищённые маршруты — обёрнуты в Auth middleware
	mux.Handle("GET /api/v1/targets", protected(http.HandlerFunc(deps.TargetHandler.List)))
	mux.Handle("POST /api/v1/targets", protected(http.HandlerFunc(deps.TargetHandler.Create)))
	mux.Handle("GET /api/v1/targets/{id}", protected(http.HandlerFunc(deps.TargetHandler.GetByID)))
	mux.Handle("PUT /api/v1/targets/{id}", protected(http.HandlerFunc(deps.TargetHandler.Update)))
	mux.Handle("DELETE /api/v1/targets/{id}", protected(http.HandlerFunc(deps.TargetHandler.Delete)))

	return mux
}
