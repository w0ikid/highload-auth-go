package bootstrap

import (
	"context"

	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
	"github.com/w0ikid/highload-auth-go/internal/usecase/account"
	"github.com/w0ikid/highload-auth-go/pkg/infra/dragonfly"
	"github.com/w0ikid/highload-auth-go/pkg/infra/postgres"
	"go.uber.org/zap"
)

type Container struct {
	logger *zap.SugaredLogger
	
	Repositories *repository.Repository
	
	AccountDomain account.AccountDomain
}

func NewContainer(
	ctx context.Context,
	pg *postgres.Postgres,
	df *dragonfly.Dragonfly,
	logger *zap.SugaredLogger,
) *Container {
	logger = logger.Named("container")

	repositories := repository.NewRepository(pg)

	baseusecase := usecase.NewBaseUsecase(repositories.ContextTransaction, logger)
	accountDomain := account.NewDomain(baseusecase, repositories.User)

	return &Container{
		logger:        logger,
		Repositories:  repositories,
		AccountDomain: accountDomain,
	}
}
