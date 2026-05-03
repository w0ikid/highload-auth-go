package rest

import (
	v1 "github.com/w0ikid/highload-auth-go/internal/transport/rest/v1"
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/account"
)

type Dependencies struct {
	AccountDeps account.HandlerDeps
}

type Handlers struct {
	V1 *v1.Handlers
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		V1: v1.NewHandlers(v1.Dependencies{
			AccountDeps: deps.AccountDeps,
		}),
	}
}
