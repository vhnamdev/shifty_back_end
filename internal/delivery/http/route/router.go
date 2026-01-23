package route

import (
	"shifty-backend/internal/delivery/http/middleware"
	"shifty-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
)

type AppHandlers struct {
}

func SetupRoutes(app *fiber.App, handlers *AppHandlers, tokenMaster *token.TokenMaster) {
	api := app.Group("/api/v1", middleware.Protected(tokenMaster))
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Shifty Backend is Ready!",
		})
	})

}
