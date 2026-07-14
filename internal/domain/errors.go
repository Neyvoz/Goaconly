package domain

import "errors"

// Доменные ошибки — не содержат деталей инфраструктуры.
// Слой usecase возвращает именно их, delivery — транслирует в HTTP-коды.
var (
	ErrTargetNotFound     = errors.New("target not found")
	ErrTargetExists       = errors.New("target with this URL already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidInterval    = errors.New("check interval must be between 1 and 1440 minutes")
	ErrTariffLimitReached = errors.New("tariff plan target limit reached")
	ErrNoCertificate      = errors.New("no TLS certificates in response")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrPasswordTooLong    = errors.New("password must not exceed 72 bytes")
)
