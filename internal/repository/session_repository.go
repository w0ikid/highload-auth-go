package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/w0ikid/highload-auth-go/pkg/infra/dragonfly"
)

type sessionRepository struct {
	df *dragonfly.Dragonfly
}

func NewSessionRepo(df *dragonfly.Dragonfly) ISessionRepo {
	return &sessionRepository{df: df}
}

func (r *sessionRepository) SetRefreshToken(ctx context.Context, token, userID string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", token)
	if err := r.df.Client.Set(ctx, key, userID, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save refresh token to dragonfly: %w", err)
	}
	return nil
}

func (r *sessionRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("refresh_token:%s", token)
	userID, err := r.df.Client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("refresh token not found or expired")
	}
	return userID, nil
}

func (r *sessionRepository) DeleteToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh_token:%s", token)
	if err := r.df.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete refresh token from dragonfly: %w", err)
	}
	return nil
}
