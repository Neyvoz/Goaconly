package usecase

import (
	"context"
	"sitepulse/internal/domain"
)

type TargetRepository interface {
	Create(ctx context.Context, target domain.Target) (domain.Target, error)
	GetByID(ctx context.Context, id int64) (domain.Target, error)
	GetAllActive(ctx context.Context) ([]domain.Target, error)
	List(ctx context.Context, userID int64, limit, offset int) ([]domain.Target, int, error)
	Update(ctx context.Context, target domain.Target) (domain.Target, error)
	Delete(ctx context.Context, id int64) error
}

type CheckResultRepository interface {
	Save(ctx context.Context, result domain.CheckResult) error                                        // ctx, не ctv
	GetLatestByTargetID(ctx context.Context, targetID int64, limit int) ([]domain.CheckResult, error) // Latest, не Lates
}
