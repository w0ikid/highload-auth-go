package dragonfly

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/w0ikid/highload-auth-go/pkg/config"
	"go.uber.org/zap"
)

type Dragonfly struct {
	Client *redis.Client
}

func New(ctx context.Context, cfg config.DragonflyConfig, logger *zap.SugaredLogger) (*Dragonfly, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to ping dragonfly: %w", err)
	}

	logger.Info("connected to dragonfly")

	return &Dragonfly{Client: client}, nil
}

func (d *Dragonfly) Close() {
	if d.Client != nil {
		d.Client.Close()
	}
}
