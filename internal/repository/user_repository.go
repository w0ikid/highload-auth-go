package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/highload-auth-go/pkg/infra/postgres"
	"github.com/w0ikid/highload-auth-go/pkg/models"
)

type userRepository struct {
	pg *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) IUserRepo {
	return &userRepository{pg: pg}
}

// db возвращает транзакцию из контекста (если она есть) 
// или обычный пул коннектов (если её нет)
func (r *userRepository) db(ctx context.Context) postgres.DBTx {
	if tx := postgres.RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.pg.Pool
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.db(ctx).QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, is_active, created_at
		FROM users
		WHERE id = $1
	`
	user := &models.User{}
	err := r.db(ctx).QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, is_active, created_at
		FROM users
		WHERE email = $1
	`
	user := &models.User{}
	err := r.db(ctx).QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}
