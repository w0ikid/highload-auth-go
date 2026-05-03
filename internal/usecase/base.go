package usecase

import (
	"github.com/w0ikid/highload-auth-go/internal/repository"
	"go.uber.org/zap"
)

type BaseUsecase struct {
	Logger *zap.SugaredLogger
	Tx     repository.IContextTransaction
}

func NewBaseUsecase(tx repository.IContextTransaction, logger *zap.SugaredLogger) BaseUsecase {
	return BaseUsecase{
		Logger: logger,
		Tx:     tx,
	}
}
