package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/highload-auth-go/pkg/infra/postgres"
	"github.com/w0ikid/highload-auth-go/pkg/models"
)

type IContextTransaction interface {
	StartTransaction(ctx context.Context) (context.Context, error)
	FinalizeTransaction(ctx context.Context, err *error) error
}

type IUserRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type Repository struct {
	ContextTransaction IContextTransaction
	User               IUserRepo
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		ContextTransaction: postgres.NewContextTransaction(pg.Pool),
		User:               NewUserRepo(pg),
	}
}
