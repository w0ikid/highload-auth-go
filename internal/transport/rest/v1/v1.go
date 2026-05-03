package v1

import (
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/account"
	"go.uber.org/zap"
)

type Dependencies struct {
	Logger      *zap.SugaredLogger
	AccountDeps account.HandlerDeps
}

type Handlers struct {
	Account account.Handler
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		Account: account.NewHandler(deps.AccountDeps),
	}
}
