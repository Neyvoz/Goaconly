package httpserver

import (
	"goaconly/internal/delivery/http/handler"
	"net/http"
)

type Dependencies struct {
	TargetHandler *handler.TargetHandler
	AuthHandler   *handler.AuthHandler
}

func NewRouter(deps Dependencies) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/targets", deps.TargetHandler.List)
	mux.HandleFunc("POST /api/v1/targets", deps.TargetHandler.Create)
	mux.HandleFunc("GET /api/v1/targets/{id}", deps.TargetHandler.GetByID)
	mux.HandleFunc("PUT /api/v1/targets/{id}", deps.TargetHandler.Update)
	mux.HandleFunc("DELETE /api/v1/targets/{id}", deps.TargetHandler.Delete)
	mux.HandleFunc("POST /api/v1/auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", deps.AuthHandler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", deps.AuthHandler.Logout)
	return mux
}
