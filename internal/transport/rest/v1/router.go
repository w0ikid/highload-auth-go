package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/accounts"
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/auth"
	"github.com/w0ikid/highload-auth-go/pkg/auth/middleware"
	"go.uber.org/zap"
)

type Router struct {
	router  fiber.Router
	handler *Handlers
}

func NewRouter(router fiber.Router, handler *Handlers) *Router {
	return &Router{
		router:  router,
		handler: handler,
	}
}

func (r *Router) SetupRoutes(logger *zap.SugaredLogger) {
	v1Router := r.router.Group("/v1")

	v1Router.Get("/ping", func(c *fiber.Ctx) error {
		logger.Info("ping received", zap.Time("time", time.Now()))
		return c.Status(200).JSON(fiber.Map{"message": "pong"})
	})

	authRouter := v1Router.Group("/auth")
	auth.NewRouter(authRouter, r.handler.Auth).SetupRoutes()

	accountsRouter := v1Router.Group("/accounts")
	accountsRouter.Use(middleware.AuthMiddleware(r.handler.JWTSecret))
	accounts.NewRouter(accountsRouter, r.handler.Accounts).SetupRoutes()
}
