package auth

import (
	"context"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
	"github.com/w0ikid/highload-auth-go/pkg/models"
)

type RegisterUsecase struct {
	usecase.BaseUsecase
	userRepo   repository.IUserRepo
	cryptoPool ICryptoPool
}

func NewRegisterUsecase(base usecase.BaseUsecase, userRepo repository.IUserRepo, cryptoPool ICryptoPool) RegisterUsecase {
	return RegisterUsecase{
		BaseUsecase: base,
		userRepo:    userRepo,
		cryptoPool:  cryptoPool,
	}
}

func (uc *RegisterUsecase) Execute(ctx context.Context, email, password string) (string, error) {
	uc.Logger.Infow("starting RegisterUsecase execution", "email", email)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = uc.Tx.FinalizeTransaction(txCtx, &err)
	}()

	// 3. Хэшируем пароль через пул воркеров
	hashedPassword, err := uc.cryptoPool.HashPassword(ctx, password)
	if err != nil {
		uc.Logger.Errorw("failed to hash password", "error", err)
		return "", err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	if err = uc.userRepo.Create(txCtx, user); err != nil {
		uc.Logger.Errorw("failed to create user", "email", email, "error", err)
		return "", err
	}

	uc.Logger.Infow("RegisterUsecase executed successfully", "user_id", user.ID)
	return user.ID.String(), nil
}
