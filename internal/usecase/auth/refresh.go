package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
	"github.com/w0ikid/highload-auth-go/pkg/crypto/jwt"
)

type RefreshUsecase struct {
	usecase.BaseUsecase
	userRepo    repository.IUserRepo
	sessionRepo repository.ISessionRepo
	jwtSecret   string
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewRefreshUsecase(
	base usecase.BaseUsecase,
	userRepo repository.IUserRepo,
	sessionRepo repository.ISessionRepo,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) RefreshUsecase {
	return RefreshUsecase{
		BaseUsecase: base,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

func (uc *RefreshUsecase) Execute(ctx context.Context, refreshToken string) (*TokenPair, error) {
	uc.Logger.Infow("starting RefreshUsecase execution")

	// 1. Проверяем валидность старого refresh токена в Dragonfly
	userIDStr, err := uc.sessionRepo.GetUserIDByToken(ctx, refreshToken)
	if err != nil {
		uc.Logger.Warnw("invalid or expired refresh token attempt")
		return nil, errors.New("unauthorized: invalid refresh token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	// 2. Получаем пользователя из БД
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.Logger.Errorw("user not found for valid session", "user_id", userID)
		return nil, errors.New("unauthorized: user not found")
	}

	// 3. Удаляем старый refresh токен из Dragonfly (Ротация)
	if err := uc.sessionRepo.DeleteToken(ctx, refreshToken); err != nil {
		// Логируем, но не прерываем выполнение, чтобы не блокировать юзера
		uc.Logger.Errorw("failed to delete old refresh token", "error", err)
	}

	// 4. Генерируем новые токены
	newAccessToken, err := jwt.GenerateToken(user.ID.String(), user.Email, uc.jwtSecret, uc.accessTTL)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 5. Сохраняем новый refresh токен
	if err := uc.sessionRepo.SetRefreshToken(ctx, newRefreshToken, user.ID.String(), uc.refreshTTL); err != nil {
		return nil, err
	}

	uc.Logger.Infow("RefreshUsecase executed successfully", "user_id", user.ID)
	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       user.ID.String(),
	}, nil
}
