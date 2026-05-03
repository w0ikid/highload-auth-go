package rest

import (
	v1 "github.com/w0ikid/highload-auth-go/internal/transport/rest/v1"
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/accounts"
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/auth"
)

type Dependencies struct {
	AuthDeps     auth.HandlerDeps
	AccountsDeps accounts.HandlerDeps
	JWTSecret    string
}

type Handlers struct {
	V1 *v1.Handlers
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		V1: v1.NewHandlers(v1.Dependencies{
			AuthDeps:     deps.AuthDeps,
			AccountsDeps: deps.AccountsDeps,
			JWTSecret:    deps.JWTSecret,
		}),
	}
}
