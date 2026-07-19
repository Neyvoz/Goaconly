package usecase

import (
	"context"
	"time"

	"goaconly/internal/domain"

	"github.com/google/uuid"

	usecaseAuth "goaconly/internal/usecase/auth"
)

type AuthUsecase interface {
	Register(ctx context.Context, rawEmail, password, companyName string) (domain.User, error)
	Login(ctx context.Context, rawEmail, password string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
}

type authUsecase struct {
	userRepo         UserRepository
	refreshTokenRepo RefreshTokenRepository
	hasher           usecaseAuth.PasswordHasher
	jwt              usecaseAuth.JWTService
	refreshTokenTTL  time.Duration
}

func (a *authUsecase) Login(ctx context.Context, rawEmail string, password string) (accessToken string, refreshToken string, err error) {
	const dummyHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	email, err := domain.NewEmail(rawEmail)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		_ = a.hasher.Compare(dummyHash, password)
		return "", "", domain.ErrInvalidCredentials
	}

	if err := a.hasher.Compare(user.PasswordHash, password); err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", "", domain.ErrInvalidCredentials
	}

	// Генерация новой пары токенов
	accessToken, err = a.jwt.GenerateAccessToken(user.ID, user.PlanID)
	if err != nil {
		return "", "", err
	}

	rawRefresh, hash, err := generateOpaqueToken()
	if err != nil {
		return "", "", err
	}

	rt := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hash,
		ExpireAt:  time.Now().Add(a.refreshTokenTTL),
		Revoked:   false,
	}
	if err := a.refreshTokenRepo.Store(ctx, rt); err != nil {
		return "", "", err
	}
	return accessToken, rawRefresh, nil
}

func (a *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)
	storedToken, err := a.refreshTokenRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil
	}
	if storedToken.Revoked {
		return nil
	}
	return a.refreshTokenRepo.Revoke(ctx, storedToken.ID)
}

func (a *authUsecase) Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	// Хэшируем токен
	hash := hashToken(refreshToken)
	// Проходимся по бд и ищем токен
	storedToken, err := a.refreshTokenRepo.GetByHash(ctx, hash)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}
	// Проверка, отозвался ли токен
	// Если токен уже отозван — просто отказываем. Более строгая защита
	// (revoke всех токенов юзера при обнаружении повторного использования
	// отозванного токена — признак возможной кражи) пока не реализована,
	// это осознанное упрощение для MVP.
	if storedToken.Revoked {
		return "", "", domain.ErrInvalidCredentials
	}
	// Проверка срока действия
	if time.Now().After(storedToken.ExpireAt) {
		return "", "", domain.ErrInvalidCredentials
	}
	// Отзываем старый токен и получаем пользователя
	if err := a.refreshTokenRepo.Revoke(ctx, storedToken.ID); err != nil {
		return "", "", err
	}
	user, err := a.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return "", "", err
	}
	// Генерация новой пары токенов как в Login
	newAccessToken, err = a.jwt.GenerateAccessToken(user.ID, user.PlanID)
	if err != nil {
		return "", "", err
	}

	rawRefresh, newHash, err := generateOpaqueToken()
	if err != nil {
		return "", "", err
	}

	rt := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: newHash,
		ExpireAt:  time.Now().Add(a.refreshTokenTTL),
		Revoked:   false,
	}
	if err := a.refreshTokenRepo.Store(ctx, rt); err != nil {
		return "", "", err
	}
	return newAccessToken, rawRefresh, nil
}

func (a *authUsecase) Register(ctx context.Context, rawEmail string, password string, companyName string) (domain.User, error) {
	email, err := domain.NewEmail(rawEmail)
	if err != nil {
		return domain.User{}, err
	}

	if len(password) > usecaseAuth.MaxPasswordBytes {
		return domain.User{}, domain.ErrPasswordTooLong
	}

	hash, err := a.hasher.Hash(password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		Email:        email,
		PasswordHash: hash,
		CompanyName:  companyName,
		PlanID:       1,
		IsActive:     true,
	}

	return a.userRepo.Create(ctx, user)
}

func NewAuthUsecase(
	userRepo UserRepository,
	refreshTokenRepo RefreshTokenRepository,
	hasher usecaseAuth.PasswordHasher,
	jwt usecaseAuth.JWTService,
	refreshTokenTTL time.Duration,
) AuthUsecase {
	return &authUsecase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		hasher:           hasher,
		jwt:              jwt,
		refreshTokenTTL:  refreshTokenTTL,
	}
}
