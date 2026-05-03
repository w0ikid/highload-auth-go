package auth

import (
	"context"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
)

type LogoutUsecase struct {
	usecase.BaseUsecase
	sessionRepo repository.ISessionRepo
}

func NewLogoutUsecase(base usecase.BaseUsecase, sessionRepo repository.ISessionRepo) LogoutUsecase {
	return LogoutUsecase{
		BaseUsecase: base,
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUsecase) Execute(ctx context.Context, refreshToken string) error {
	uc.Logger.Infow("starting LogoutUsecase execution")
	
	if refreshToken == "" {
		return nil // Ничего не делаем, если токена нет
	}

	return uc.sessionRepo.DeleteToken(ctx, refreshToken)
}
