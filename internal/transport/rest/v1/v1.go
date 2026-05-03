package v1

import (
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/accounts"
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/auth"
	"go.uber.org/zap"
)

type Dependencies struct {
	Logger       *zap.SugaredLogger
	AuthDeps     auth.HandlerDeps
	AccountsDeps accounts.HandlerDeps
	JWTSecret    string
}

type Handlers struct {
	Auth      auth.Handler
	Accounts  accounts.Handler
	JWTSecret string
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		Auth:      auth.NewHandler(deps.AuthDeps),
		Accounts:  accounts.NewHandler(deps.AccountsDeps),
		JWTSecret: deps.JWTSecret,
	}
}
