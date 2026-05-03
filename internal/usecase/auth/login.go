package auth

import (
	"context"
	"errors"
	"time"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
	"github.com/w0ikid/highload-auth-go/pkg/crypto/jwt"
)

type LoginUsecase struct {
	usecase.BaseUsecase
	userRepo    repository.IUserRepo
	sessionRepo repository.ISessionRepo
	jwtSecret   string
	accessTTL   time.Duration
	refreshTTL  time.Duration
	cryptoPool  ICryptoPool
}

func NewLoginUsecase(
	base usecase.BaseUsecase,
	userRepo repository.IUserRepo,
	sessionRepo repository.ISessionRepo,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	cryptoPool ICryptoPool,
) LoginUsecase {
	return LoginUsecase{
		BaseUsecase: base,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
		cryptoPool:  cryptoPool,
	}
}

type TokenPair struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (uc *LoginUsecase) Execute(ctx context.Context, email, password string) (*TokenPair, error) {
	uc.Logger.Infow("starting LoginUsecase execution", "email", email)

	// 1. Получаем пользователя по email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		uc.Logger.Warnw("user not found", "email", email)
		return nil, errors.New("invalid email or password") // Не выдаем, что именно не так
	}

	// 2. Сравниваем хэши через пул воркеров
	match, err := uc.cryptoPool.ComparePassword(ctx, password, user.PasswordHash)
	if err != nil || !match {
		uc.Logger.Warnw("invalid password attempt", "email", email)
		return nil, errors.New("invalid email or password")
	}

	// 3. Генерируем Access JWT токен
	accessToken, err := jwt.GenerateToken(user.ID.String(), user.Email, uc.jwtSecret, uc.accessTTL)
	if err != nil {
		uc.Logger.Errorw("failed to generate access token", "user_id", user.ID, "error", err)
		return nil, err
	}

	// 4. Генерируем Opaque Refresh токен
	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		uc.Logger.Errorw("failed to generate refresh token", "user_id", user.ID, "error", err)
		return nil, err
	}

	// 5. Сохраняем Refresh токен в сессию (Dragonfly)
	if err := uc.sessionRepo.SetRefreshToken(ctx, refreshToken, user.ID.String(), uc.refreshTTL); err != nil {
		uc.Logger.Errorw("failed to save session", "user_id", user.ID, "error", err)
		return nil, err
	}

	uc.Logger.Infow("LoginUsecase executed successfully", "email", email, "user_id", user.ID)
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID.String(),
	}, nil
}
