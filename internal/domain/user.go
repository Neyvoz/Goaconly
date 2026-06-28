package domain

import "time"

// UserRole — роль пользователя в системе.
// Строка, не iota — по той же причине что и CheckStatus.
type UserRole string

const (
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
	RoleViewer UserRole = "viewer"
)

// User — пользователь системы SitePulse.
// Пароль хранится только в виде хеша — никогда не храним plaintext.
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         UserRole
	PlanID       int64
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
