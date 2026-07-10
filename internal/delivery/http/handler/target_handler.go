package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sitepulse/internal/domain"
	"sitepulse/internal/usecase"
)

// TargetHandler — HTTP-обёртка над бизнес-логикой TargetUsecase.
// Не содержит бизнес-правил — только парсинг запроса, вызов usecase,
// маппинг ошибок в HTTP-статусы и сериализацию ответа.
type TargetHandler struct {
	uc usecase.TargetUsecase
}

func NewTargetHandler(uc usecase.TargetUsecase) *TargetHandler {
	return &TargetHandler{uc: uc}
}

func (h *TargetHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateTargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, errors.New("invalid JSON body"))
		return
	}

	// Валидация ФОРМАТА (не бизнес-правил) — задача хендлера
	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		responseError(w, http.StatusBadRequest, errors.New("url must be a valid http/https URL"))
		return
	}

	userID := getUserIDFromContext(r.Context())

	target, err := h.uc.Create(r.Context(), userID, req.URL, req.KeywordToFind, req.CheckIntervalMin)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidInterval):
			responseError(w, http.StatusUnprocessableEntity, err)
		case errors.Is(err, domain.ErrTargetExists):
			responseError(w, http.StatusConflict, err)
		default:
			responseError(w, http.StatusInternalServerError, errors.New("internal server error"))
		}
		return
	}

	responseJSON(w, http.StatusCreated, toTargetResponse(target))
}

// TODO: временная заглушка до Недели 6 (JWT). Возвращает фиксированный userID
// для локальной разработки без авторизации.
func getUserIDFromContext(ctx context.Context) int64 {
	_ = ctx
	return 1
}
