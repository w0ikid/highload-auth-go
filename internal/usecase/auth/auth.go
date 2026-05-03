package auth

import (
	"context"
	"time"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
)

type ICryptoPool interface {
	HashPassword(ctx context.Context, password string) (string, error)
	ComparePassword(ctx context.Context, password, hash string) (bool, error)
}

type AuthDomain struct {
	RegisterUsecase RegisterUsecase
	LoginUsecase    LoginUsecase
	RefreshUsecase  RefreshUsecase
	LogoutUsecase   LogoutUsecase
}

func NewAuthDomain(
	baseusecase usecase.BaseUsecase,
	userRepo repository.IUserRepo,
	sessionRepo repository.ISessionRepo,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	cryptoPool ICryptoPool,
) AuthDomain {
	baseusecase.Logger = baseusecase.Logger.Named("auth_domain")
	return AuthDomain{
		RegisterUsecase: NewRegisterUsecase(baseusecase, userRepo, cryptoPool),
		LoginUsecase:    NewLoginUsecase(baseusecase, userRepo, sessionRepo, jwtSecret, accessTTL, refreshTTL, cryptoPool),
		RefreshUsecase:  NewRefreshUsecase(baseusecase, userRepo, sessionRepo, jwtSecret, accessTTL, refreshTTL),
		LogoutUsecase:   NewLogoutUsecase(baseusecase, sessionRepo),
	}
}
