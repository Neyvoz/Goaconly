package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserRole — роль пользователя в системе.
// Строка, не iota — по той же причине что и CheckStatus.
type UserRole string

// RBAC — вернуть Role в структуру User вместе с сущностью Company.
const (
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
	RoleViewer UserRole = "viewer"
)

// User — пользователь системы SitePulse.
// Пароль хранится только в виде хеша — никогда не храним plaintext.
type User struct {
	ID           uuid.UUID
	Email        Email
	PasswordHash string
	CompanyName  string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Email — Value Object. Гарантирует, что в системе не существует
// User с невалидным email, независимо от точки создания.
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail нормализует email (lowercase, trim).
func NewEmail(raw string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return Email{}, fmt.Errorf("%w: email is empty", ErrInvalidEmail)
	}
	if !emailRegex.MatchString(normalized) {
		return Email{}, fmt.Errorf("%w: %q", ErrInvalidEmail, raw)
	}
	return Email{value: normalized}, nil
}

// String возвращает нормализованное строковое представление.
func (e Email) String() string {
	return e.value
}
