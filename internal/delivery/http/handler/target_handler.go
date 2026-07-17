package handler

import (
	"context"
	"encoding/json"
	"errors"
	"goaconly/internal/domain"
	"goaconly/internal/usecase"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// TargetHandler — HTTP-обёртка над бизнес-логикой TargetUsecase.
// Не содержит бизнес-правил — только парсинг запроса, вызов usecase,
// маппинг ошибок в HTTP-статусы и сериализацию ответа.
type TargetHandler struct {
	uc usecase.TargetUsecase
}

// NewTargetHandler — конструктор, через который main.go внедряет зависимость
func NewTargetHandler(uc usecase.TargetUsecase) *TargetHandler {
	return &TargetHandler{uc: uc}
}

// Create — POST /api/v1/targets
// Создаёт новую цель мониторинга для текущего пользователя.
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

// GetByID — GET /api/v1/targets/{id}
// Возвращает одну цель мониторинга по её ID, если она принадлежит
// текущему пользователю.
func (h *TargetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	targetID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		responseError(w, http.StatusBadRequest, errors.New("invalid target id"))
		return
	}
	userID := getUserIDFromContext(r.Context())

	target, err := h.uc.GetByID(r.Context(), userID, targetID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTargetNotFound):
			responseError(w, http.StatusNotFound, err)
		default:
			responseError(w, http.StatusInternalServerError, errors.New("internal server error"))
		}
		return
	}
	responseJSON(w, http.StatusOK, toTargetResponse(target))
}

// List — GET /api/v1/targets?limit=20&offset=0
// Возвращает страницу целей мониторинга текущего пользователя с пагинацией.
func (h *TargetHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 20)
	offset := parseIntQuery(r, "offset", 0)

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	userID := getUserIDFromContext(r.Context())
	targets, total, err := h.uc.List(r.Context(), userID, limit, offset)
	if err != nil {
		responseError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	items := make([]TargetResponse, len(targets))
	for i, t := range targets {
		items[i] = toTargetResponse(t)
	}
	responseJSON(w, http.StatusOK, ListResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// Update — PUT /api/v1/targets/{id}
// Обновляет редактируемые поля цели мониторинга (URL, ключевое слово,
// интервал проверки), если она принадлежит текущему пользователю.
func (h *TargetHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	targetID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		responseError(w, http.StatusBadRequest, errors.New("invalid target id"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateTargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, errors.New("invalid JSON body"))
		return
	}
	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		responseError(w, http.StatusBadRequest, errors.New("url must be a valid http/https URL"))
		return
	}
	userID := getUserIDFromContext(r.Context())

	target, err := h.uc.Update(r.Context(), userID, targetID, req.URL, req.KeywordToFind, req.CheckIntervalMin)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTargetNotFound):
			responseError(w, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrInvalidInterval):
			responseError(w, http.StatusUnprocessableEntity, err)
		default:
			responseError(w, http.StatusInternalServerError, errors.New("internal server error"))
		}
		return
	}
	responseJSON(w, http.StatusOK, toTargetResponse(target))
}

// Delete — DELETE /api/v1/targets/{id}
// Удаляет цель мониторинга по ID, если она принадлежит текущему пользователю.
func (h *TargetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	targetID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		responseError(w, http.StatusBadRequest, errors.New("invalid target id"))
		return
	}
	userID := getUserIDFromContext(r.Context())
	err = h.uc.Delete(r.Context(), userID, targetID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTargetNotFound):
			responseError(w, http.StatusNotFound, err)
		default:
			responseError(w, http.StatusInternalServerError, errors.New("internal server error"))
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// TODO: временная заглушка до Недели 6 (JWT). Возвращает фиксированный userID
// для локальной разработки без авторизации.
func getUserIDFromContext(ctx context.Context) uuid.UUID {
	_ = ctx
	id, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	return id
}
