package accounts

import (
	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
)

type AccountsDomain struct {
	GetProfileUsecase GetProfileUsecase
}

func NewAccountsDomain(baseusecase usecase.BaseUsecase, userRepo repository.IUserRepo) AccountsDomain {
	baseusecase.Logger = baseusecase.Logger.Named("accounts_domain")
	return AccountsDomain{
		GetProfileUsecase: NewGetProfileUsecase(baseusecase, userRepo),
	}
}
