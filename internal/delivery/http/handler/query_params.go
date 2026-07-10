package handler

import (
	"net/http"
	"strconv"
)

// parseIntQuery читает int-параметр из query-строки с дефолтным значением.
func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}
	return value
}
