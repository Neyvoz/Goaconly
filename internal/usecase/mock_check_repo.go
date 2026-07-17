package usecase

import (
	"context"
	"goaconly/internal/domain"
	"log"
)

type MockCheckRepository struct{}

func (r *MockCheckRepository) Save(ctx context.Context, result domain.CheckResult) error {
	log.Printf("[MockRepo] Saved result: targetID=%d status=%d responseTime=%dms",
		result.TargetID,
		result.StatusCode,
		result.ResponseTimeMs,
	)
	return nil
}

func (r *MockCheckRepository) GetLatestByTargetID(ctx context.Context, targetID int64, limit int) ([]domain.CheckResult, error) {
	return nil, nil
}

var _ CheckResultRepository = (*MockCheckRepository)(nil)
