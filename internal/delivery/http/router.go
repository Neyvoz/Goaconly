package httpserver

import "net/http"

type Depenencies struct {
}

func NewRouter(deps Depenencies) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/targets", placeholderHandler)
	mux.HandleFunc("POST /api/v1/targets", placeholderHandler)
	mux.HandleFunc("GET /api/v1/targets/{id}", placeholderHandler)
	mux.HandleFunc("PUT /api/v1/targets/{id}", placeholderHandler)
	mux.HandleFunc("DELETE /api/v1/targets/{id}", placeholderHandler)
	return mux
}
func placeholderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
