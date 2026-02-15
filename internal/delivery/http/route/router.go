package route

import (
	"shifty-backend/internal/delivery/http/handler"
	"shifty-backend/internal/delivery/http/middleware"
	"shifty-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
)

type AppHandlers struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
}

func SetupRoutes(app *fiber.App, h *AppHandlers, tokenMaster *token.TokenMaster) {
	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Shifty Backend is Ready!",
		})
	})

	// -----------------------------AUTH GROUP----------------------------------
	auth := api.Group("/auth")
	auth.Post("/register", h.AuthHandler.RegisterLocal)
	auth.Post("/login", h.AuthHandler.LoginLocal)
	auth.Post("/send-otp", h.AuthHandler.SendOTP)

	protected := api.Group("/", middleware.Protected(tokenMaster))
	protected.Get("/pofile", h.UserHandler.UpdateAvatar)
}
