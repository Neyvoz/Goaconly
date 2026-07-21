package handler

import (
	"encoding/json"
	"errors"
	"goaconly/internal/domain"
	"goaconly/internal/usecase"
	"net/http"
	"time"
)

type AuthHandler struct {
	uc usecase.AuthUsecase
}

func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, errors.New("invalid JSON body"))
		return
	}

	user, err := h.uc.Register(r.Context(), req.Email, req.Password, req.CompanyName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			responseError(w, http.StatusConflict, err)
		case errors.Is(err, domain.ErrInvalidEmail):
			responseError(w, http.StatusUnprocessableEntity, err)
		case errors.Is(err, domain.ErrPasswordTooLong):
			responseError(w, http.StatusUnprocessableEntity, err)
		case errors.Is(err, domain.ErrPasswordTooShort):
			responseError(w, http.StatusUnprocessableEntity, err)
		default:
			responseError(w, http.StatusInternalServerError, errors.New("internal server error"))
		}
		return
	}

	responseJSON(w, http.StatusCreated, RegisterResponse{
		ID:          user.ID.String(),
		Email:       user.Email.String(),
		CompanyName: user.CompanyName,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, errors.New("invalid JSON body"))
		return
	}

	accessToken, refreshToken, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			responseError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		default:
			responseError(w, http.StatusInternalServerError, errors.New("internal server error"))
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
	})

	responseJSON(w, http.StatusOK, LoginResponse{AccessToken: accessToken})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		responseError(w, http.StatusUnauthorized, errors.New("refresh token cookie not found"))
		return
	}

	newAccessToken, newRefreshToken, err := h.uc.Refresh(r.Context(), cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "refresh_token",
			Value:  "",
			Path:   "/api/v1/auth",
			MaxAge: -1,
		})
		responseError(w, http.StatusUnauthorized, errors.New("invalid or expired refresh token"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
	})

	responseJSON(w, http.StatusOK, RefreshResponse{AccessToken: newAccessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.uc.Logout(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		Value:  "",
		Path:   "/api/v1/auth",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
