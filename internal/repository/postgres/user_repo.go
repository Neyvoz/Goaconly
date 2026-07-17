package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sitepulse/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, t domain.User) (domain.User, error) {
	const q = `
		INSERT INTO users (email, password_hash, company_name, plan_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, q,
		t.Email.String(), t.PasswordHash, t.CompanyName, t.PlanID, t.IsActive,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.User{}, domain.ErrUserAlreadyExists
		}
		return domain.User{}, err
	}
	return t, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	const q = `
		SELECT id, email, password_hash, company_name, plan_id, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var t domain.User
	var rawEmail string

	err := r.db.QueryRowContext(ctx, q, email.String()).Scan(
		&t.ID, &rawEmail, &t.PasswordHash, &t.CompanyName,
		&t.PlanID, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	t.Email, err = domain.NewEmail(rawEmail)
	if err != nil {
		return domain.User{}, err
	}
	return t, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	const q = `
		SELECT id, email, password_hash, company_name, plan_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var t domain.User
	var email string

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&t.ID, &email, &t.PasswordHash, &t.CompanyName,
		&t.PlanID, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	t.Email, err = domain.NewEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	return t, nil
}
