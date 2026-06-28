package postgres

import (
	"context"
	"database/sql"
	"sitepulse/internal/domain"
)

type checkResultRepo struct {
	db *sql.DB
}

func NewCheckResultRepo(db *sql.DB) *checkResultRepo {
	return &checkResultRepo{db: db}
}

func (r *checkResultRepo) Save(ctx context.Context, result domain.CheckResult) error {
	const q = `
        INSERT INTO check_logs
            (target_id, checked_at, status_code, response_time_ms, is_up, ssl_expires_at, error_message)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(ctx, q,
		result.TargetID,
		result.CheckedAt,
		result.StatusCode,
		result.ResponseTimeMs,
		result.IsUp,
		result.SSLExpiresAt,
		result.ErrorMessage,
	)
	return err
}

func (r *checkResultRepo) GetLatestByTargetID(ctx context.Context, targetID int64, limit int) ([]domain.CheckResult, error) {
	const q = `
        SELECT target_id, checked_at, status_code, response_time_ms,
               is_up, ssl_expires_at, error_message
        FROM check_logs
        WHERE target_id = $1
        ORDER BY checked_at DESC
        LIMIT $2
    `
	rows, err := r.db.QueryContext(ctx, q, targetID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.CheckResult
	for rows.Next() {
		var cr domain.CheckResult
		if err := rows.Scan(
			&cr.TargetID,
			&cr.CheckedAt,
			&cr.StatusCode,
			&cr.ResponseTimeMs,
			&cr.IsUp,
			&cr.SSLExpiresAt,
			&cr.ErrorMessage,
		); err != nil {
			return nil, err
		}
		results = append(results, cr)
	}
	return results, rows.Err()
}
