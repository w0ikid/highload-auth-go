package accounts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/highload-auth-go/internal/usecase/accounts"
)

type HandlerDeps struct {
	AccountsDomain accounts.AccountsDomain
}

type Handler struct {
	accountsDomain accounts.AccountsDomain
}

func NewHandler(deps HandlerDeps) Handler {
	return Handler{
		accountsDomain: deps.AccountsDomain,
	}
}

func (h *Handler) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	profile, err := h.accountsDomain.GetProfileUsecase.Execute(c.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}
