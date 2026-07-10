package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse — единый формат ошибок API.
// Используется во всех хендлерах проекта, не только Target.
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondJSON — централизованная запись JSON-ответа.
// Гарантирует правильный порядок: сначала заголовки, потом статус, потом тело.
func responseJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondError — обёртка над respondJSON для ошибок.
// ВАЖНО: сюда нельзя передавать сырые внутренние ошибки (например, от sql.DB) —
// они могут содержать детали инфраструктуры. Для непредвиденных ошибок
// хендлер должен подставлять generic-сообщение, а исходную ошибку логировать отдельно.
func responseError(w http.ResponseWriter, status int, err error) {
	responseJSON(w, status, ErrorResponse{Error: err.Error()})
}
