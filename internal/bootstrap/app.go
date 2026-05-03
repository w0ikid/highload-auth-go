package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/w0ikid/highload-auth-go/pkg/config"
	"github.com/w0ikid/highload-auth-go/pkg/infra/dragonfly"
	"github.com/w0ikid/highload-auth-go/pkg/infra/postgres"

	"github.com/w0ikid/highload-auth-go/internal/transport/rest"
	accountsHandler "github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/accounts"
	authHandler "github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/auth"
)

type App struct {
	fapp      *fiber.App
	addr      string
	container *Container
	logger    *zap.SugaredLogger
	pg        *postgres.Postgres
	dragonfly *dragonfly.Dragonfly
	cancel    context.CancelFunc
}

func NewApp(ctx context.Context, cfg config.Config, logger *zap.SugaredLogger) (*App, error) {
	appLogger := logger.Named("app")

	// postgres
	pg, err := postgres.New(ctx, cfg.Postgres, appLogger)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	// dragonfly
	df, err := dragonfly.New(ctx, cfg.Dragonfly, appLogger)
	if err != nil {
		return nil, fmt.Errorf("connect dragonfly: %w", err)
	}

	// DI Container
	cont := NewContainer(ctx, cfg, pg, df, appLogger)

	fapp := fiber.New(fiber.Config{
		AppName:      "highload-auth-go",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	fapp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			appLogger.Errorw("panic recovered", "error", e)
		},
	}))

	// Handlers DI
	h := rest.NewHandlers(rest.Dependencies{
		AuthDeps: authHandler.HandlerDeps{
			AuthDomain: cont.AuthDomain,
		},
		AccountsDeps: accountsHandler.HandlerDeps{
			AccountsDomain: cont.AccountsDomain,
		},
		JWTSecret: cfg.JWT.Secret,
	})

	// Fiber router
	router := rest.NewRouter(fapp, h)
	router.SetupRoutes(appLogger)

	_, cancel := context.WithCancel(ctx)

	return &App{
		fapp:      fapp,
		addr:      ":" + cfg.HTTP.Port,
		container: cont,
		logger:    appLogger,
		pg:        pg,
		dragonfly: df,
		cancel:    cancel,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info("starting HTTP server", zap.String("addr", a.addr))
	if err := a.fapp.Listen(a.addr); err != nil {
		return fmt.Errorf("fiber server: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.cancel()

	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var errOccurred bool

	if err := a.fapp.ShutdownWithContext(shutdownCtx); err != nil {
		a.logger.Error("fiber shutdown failed", zap.Error(err))
		errOccurred = true
	} else {
		a.logger.Info("fiber server stopped gracefully")
	}

	a.pg.Close()
	a.logger.Info("postgres connection closed")

	a.dragonfly.Close()
	a.logger.Info("dragonfly connection closed")

	if errOccurred {
		return fmt.Errorf("some resources failed to close, check logs")
	}

	a.logger.Info("app stopped gracefully")
	return nil
}
