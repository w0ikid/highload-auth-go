package account

import (
	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
)

type AccountDomain struct {
	RegisterUsecase RegisterUsecase
	LoginUsecase    LoginUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, userRepo repository.IUserRepo) AccountDomain {
	baseusecase.Logger = baseusecase.Logger.Named("account_domain")
	return AccountDomain{
		RegisterUsecase: NewRegisterUsecase(baseusecase, userRepo),
		LoginUsecase:    NewLoginUsecase(baseusecase, userRepo),
	}
}
