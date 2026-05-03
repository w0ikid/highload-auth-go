package bootstrap

import (
	"context"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
	"github.com/w0ikid/highload-auth-go/internal/usecase/accounts"
	"github.com/w0ikid/highload-auth-go/internal/usecase/auth"
	"github.com/w0ikid/highload-auth-go/pkg/config"
	"github.com/w0ikid/highload-auth-go/pkg/infra/dragonfly"
	"github.com/w0ikid/highload-auth-go/pkg/infra/postgres"
	"go.uber.org/zap"
)

type Container struct {
	logger *zap.SugaredLogger
	
	Repositories *repository.Repository
	
	AuthDomain     auth.AuthDomain
	AccountsDomain accounts.AccountsDomain
}

func NewContainer(
	ctx context.Context,
	cfg config.Config,
	pg *postgres.Postgres,
	df *dragonfly.Dragonfly,
	logger *zap.SugaredLogger,
) *Container {
	logger = logger.Named("container")

	repositories := repository.NewRepository(pg, df)

	baseusecase := usecase.NewBaseUsecase(repositories.ContextTransaction, logger)
	authDomain := auth.NewAuthDomain(baseusecase, repositories.User, repositories.Session, cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	accountsDomain := accounts.NewAccountsDomain(baseusecase, repositories.User)

	return &Container{
		logger:         logger,
		Repositories:   repositories,
		AuthDomain:     authDomain,
		AccountsDomain: accountsDomain,
	}
}
