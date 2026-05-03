package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/highload-auth-go/internal/transport/rest/v1/account"
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

	accountRouter := v1Router.Group("/account")
	account.NewRouter(accountRouter, r.handler.Account).SetupRoutes()
}
