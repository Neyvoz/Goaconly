package postgres

import (
	"context"
	"database/sql"
	"errors"
	"goaconly/internal/domain"

	"github.com/google/uuid"
)

type refreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) *refreshTokenRepo {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Store(ctx context.Context, t domain.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, q, t.ID, t.UserID, t.TokenHash, t.ExpireAt, t.Revoked)
	return err
}

func (r *refreshTokenRepo) GetByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var t domain.RefreshToken
	err := r.db.QueryRowContext(ctx, q, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpireAt, &t.Revoked, &t.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.RefreshToken{}, domain.ErrRefreshTokenNotFound
	}
	if err != nil {
		return domain.RefreshToken{}, err
	}
	return t, nil
}

func (r *refreshTokenRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE refresh_tokens SET revoked = true WHERE id = $1`
	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return domain.ErrRefreshTokenNotFound
	}
	return nil
}

func (r *refreshTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1 AND revoked = false`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}
