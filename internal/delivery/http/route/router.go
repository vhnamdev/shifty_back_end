package route

import (
	"github.com/gofiber/fiber/v2"
)

type AppHandlers struct {
}

func SetupRoutes(app *fiber.App, handlers *AppHandlers) {
	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Shifty Backend is Ready!",
		})
	})

}
