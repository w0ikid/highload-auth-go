package rest

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/w0ikid/highload-auth-go/internal/transport/rest/v1"
	"go.uber.org/zap"
)

type Router struct {
	fapp    *fiber.App
	handler *Handlers
}

func NewRouter(fapp *fiber.App, handler *Handlers) *Router {
	return &Router{
		fapp:    fapp,
		handler: handler,
	}
}

func (r *Router) SetupRoutes(logger *zap.SugaredLogger) {
	apiRouter := r.fapp.Group("/api")
	
	v1Router := v1.NewRouter(apiRouter, r.handler.V1)
	v1Router.SetupRoutes(logger)
}
