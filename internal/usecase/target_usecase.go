package usecase

import (
	"context"
	"goaconly/internal/domain"

	"github.com/google/uuid"
)

// Методы Create, GetByID, List, Update, Delete
type TargetUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, url string, keyword string, intervalMinutes int) (domain.Target, error)
	GetByID(ctx context.Context, userID uuid.UUID, targetID int64) (domain.Target, error)
	List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Target, int, error)
	Update(ctx context.Context, userID uuid.UUID, targetID int64, url string, keyword string, intervalMinutes int) (domain.Target, error)
	Delete(ctx context.Context, userID uuid.UUID, targetID int64) error
}

// targetUsecase — приватная реализация интерфейса TargetUsecase.
type targetUsecase struct {
	repo TargetRepository
}

func (u *targetUsecase) GetByID(ctx context.Context, userID uuid.UUID, targetID int64) (domain.Target, error) {
	target, err := u.repo.GetByID(ctx, targetID)
	if err != nil {
		return domain.Target{}, err
	}

	// Бизнес-правило владения: пользователь не может получить чужую цель
	// мониторинга по ID. Это не забота репозитория (он просто читает по id) —
	// это забота usecase, потому что здесь мы знаем "кто спрашивает".
	if target.UserID != userID {
		return domain.Target{}, domain.ErrTargetNotFound
	}

	return target, nil
}

func (u *targetUsecase) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Target, int, error) {
	return u.repo.List(ctx, userID, limit, offset)
}

func (u *targetUsecase) Update(ctx context.Context, userID uuid.UUID, targetID int64, url string, keyword string, intervalMinutes int) (domain.Target, error) {
	if intervalMinutes < 1 || intervalMinutes > 1440 {
		return domain.Target{}, domain.ErrInvalidInterval
	}

	existing, err := u.repo.GetByID(ctx, targetID)
	if err != nil {
		return domain.Target{}, err
	}
	if existing.UserID != userID {
		return domain.Target{}, domain.ErrTargetNotFound
	}

	existing.URL = url
	existing.KeywordToFind = keyword
	existing.CheckInterval = domain.CheckInterval(intervalMinutes)

	return u.repo.Update(ctx, existing)
}

func (u *targetUsecase) Delete(ctx context.Context, userID uuid.UUID, targetID int64) error {
	existing, err := u.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return domain.ErrTargetNotFound
	}

	return u.repo.Delete(ctx, targetID)
}

func NewTargetUsecase(repo TargetRepository) TargetUsecase {
	return &targetUsecase{repo: repo}
}

func (u *targetUsecase) Create(ctx context.Context, userID uuid.UUID, url string, keyword string, intervalMinutes int) (domain.Target, error) {
	if intervalMinutes < 1 || intervalMinutes > 1440 {
		return domain.Target{}, domain.ErrInvalidInterval
	}
	target := domain.Target{
		UserID:        userID,
		URL:           url,
		KeywordToFind: keyword,
		CheckInterval: domain.CheckInterval(intervalMinutes),
		IsActive:      true,
	}
	created, err := u.repo.Create(ctx, target)
	if err != nil {
		return domain.Target{}, err
	}
	return created, nil
}
