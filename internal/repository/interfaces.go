package repository

import (
	"context"
	"sitepulse/internal/domain"
)

type TargetRepository interface {
	GetAllActive(ctx context.Context) ([]domain.Target, error)
}

type CheckResultRepository interface {
	Save(ctx context.Context, result domain.CheckResult) error                                        // ctx, не ctv
	GetLatestByTargetID(ctx context.Context, targetID int64, limit int) ([]domain.CheckResult, error) // Latest, не Lates
}
