package httpserver

import (
	"goaconly/internal/delivery/http/handler"
	"net/http"
)

type Dependencies struct {
	TargetHandler *handler.TargetHandler
}

func NewRouter(deps Dependencies) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/targets", deps.TargetHandler.List)
	mux.HandleFunc("POST /api/v1/targets", deps.TargetHandler.Create)
	mux.HandleFunc("GET /api/v1/targets/{id}", deps.TargetHandler.GetByID)
	mux.HandleFunc("PUT /api/v1/targets/{id}", deps.TargetHandler.Update)
	mux.HandleFunc("DELETE /api/v1/targets/{id}", deps.TargetHandler.Delete)
	return mux
}
