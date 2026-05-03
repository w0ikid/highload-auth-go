package account

import (
	"context"

	"errors"
	"time"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
	"github.com/w0ikid/highload-auth-go/pkg/crypto/hash"
	"github.com/w0ikid/highload-auth-go/pkg/crypto/jwt"
)

type LoginUsecase struct {
	usecase.BaseUsecase
	userRepo repository.IUserRepo
}

func NewLoginUsecase(base usecase.BaseUsecase, userRepo repository.IUserRepo) LoginUsecase {
	return LoginUsecase{
		BaseUsecase: base,
		userRepo:    userRepo,
	}
}

func (uc *LoginUsecase) Execute(ctx context.Context, email, password string) (string, error) {
	uc.Logger.Infow("starting LoginUsecase execution", "email", email)

	// 1. Получаем пользователя по email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		uc.Logger.Warnw("user not found", "email", email)
		return "", errors.New("invalid email or password") // Не выдаем, что именно не так (для безопасности)
	}

	// 2. Сравниваем хэши
	match, err := hash.ComparePassword(password, user.PasswordHash)
	if err != nil || !match {
		uc.Logger.Warnw("invalid password attempt", "email", email)
		return "", errors.New("invalid email or password")
	}

	// 3. Генерируем JWT токен
	// TODO: Секретный ключ нужно брать из конфига (config.go)
	secretKey := "super-secret-key-replace-me"
	token, err := jwt.GenerateToken(user.ID.String(), user.Email, secretKey, 24*time.Hour)
	if err != nil {
		uc.Logger.Errorw("failed to generate token", "user_id", user.ID, "error", err)
		return "", err
	}

	uc.Logger.Infow("LoginUsecase executed successfully", "email", email, "user_id", user.ID)
	return token, nil
}
