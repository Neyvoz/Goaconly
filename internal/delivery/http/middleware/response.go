package middleware

import (
	"encoding/json"
	"net/http"
)

// responseError — локальная копия хелпера из пакета handler.
func responseError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
