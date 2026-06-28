package postgres

import (
	"context"
	"database/sql"
	"sitepulse/internal/domain"
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
			&t.ID,
			&t.UserID,
			&t.URL,
			&t.CheckInterval,
			&t.KeywordToFind,
			&t.IsActive,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, rows.Err()
}
