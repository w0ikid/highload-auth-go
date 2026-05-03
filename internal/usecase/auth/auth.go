package auth

import (
	"time"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
)

type AuthDomain struct {
	RegisterUsecase RegisterUsecase
	LoginUsecase    LoginUsecase
	RefreshUsecase  RefreshUsecase
}

func NewAuthDomain(
	baseusecase usecase.BaseUsecase,
	userRepo repository.IUserRepo,
	sessionRepo repository.ISessionRepo,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) AuthDomain {
	baseusecase.Logger = baseusecase.Logger.Named("auth_domain")
	return AuthDomain{
		RegisterUsecase: NewRegisterUsecase(baseusecase, userRepo),
		LoginUsecase:    NewLoginUsecase(baseusecase, userRepo, sessionRepo, jwtSecret, accessTTL, refreshTTL),
		RefreshUsecase:  NewRefreshUsecase(baseusecase, userRepo, sessionRepo, jwtSecret, accessTTL, refreshTTL),
	}
}
