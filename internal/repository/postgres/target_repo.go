package postgres

import (
	"context"
	"database/sql"
	"errors"

	"goaconly/internal/domain"

	"github.com/google/uuid"
)

type targetRepo struct {
	db *sql.DB
}

func NewTargetRepo(db *sql.DB) *targetRepo {
	return &targetRepo{db: db}
}

func (r *targetRepo) GetAllActive(ctx context.Context) ([]domain.Target, error) {
	const q = `
		SELECT id, user_id, url, check_interval, keyword_to_find, is_active, created_at, updated_at
		FROM targets
		WHERE is_active = true
	`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var t domain.Target
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.URL, &t.CheckInterval,
			&t.KeywordToFind, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, rows.Err()
}

func (r *targetRepo) Create(ctx context.Context, t domain.Target) (domain.Target, error) {
	const q = `
		INSERT INTO targets (user_id, url, check_interval, keyword_to_find, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, q,
		t.UserID, t.URL, t.CheckInterval, t.KeywordToFind, t.IsActive,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return domain.Target{}, err
	}
	return t, nil
}

func (r *targetRepo) GetByID(ctx context.Context, id int64) (domain.Target, error) {
	const q = `
		SELECT id, user_id, url, check_interval, keyword_to_find, is_active, created_at, updated_at
		FROM targets
		WHERE id = $1
	`
	var t domain.Target
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&t.ID, &t.UserID, &t.URL, &t.CheckInterval,
		&t.KeywordToFind, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Target{}, domain.ErrTargetNotFound
	}
	if err != nil {
		return domain.Target{}, err
	}
	return t, nil
}

func (r *targetRepo) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Target, int, error) {
	const countQ = `SELECT COUNT(*) FROM targets WHERE user_id = $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQ, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const listQ = `
		SELECT id, user_id, url, check_interval, keyword_to_find, is_active, created_at, updated_at
		FROM targets
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, listQ, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var t domain.Target
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.URL, &t.CheckInterval,
			&t.KeywordToFind, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		targets = append(targets, t)
	}
	return targets, total, rows.Err()
}

func (r *targetRepo) Update(ctx context.Context, t domain.Target) (domain.Target, error) {
	const q = `
		UPDATE targets
		SET url = $1, check_interval = $2, keyword_to_find = $3, updated_at = now()
		WHERE id = $4
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, q,
		t.URL, t.CheckInterval, t.KeywordToFind, t.ID,
	).Scan(&t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Target{}, domain.ErrTargetNotFound
	}
	if err != nil {
		return domain.Target{}, err
	}
	return t, nil
}

func (r *targetRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM targets WHERE id = $1`
	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrTargetNotFound
	}
	return nil
}
