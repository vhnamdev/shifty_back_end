package route

import (
	"context"
	"shifty-backend/graph"
	"shifty-backend/internal/delivery/http/handler"
	"shifty-backend/internal/delivery/http/middleware"
	"shifty-backend/pkg/token"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type AppHandlers struct {
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
}

func SetupRoutes(app *fiber.App, h *AppHandlers, tokenMaster *token.TokenMaster, gqlResolver *graph.Resolver) {
	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Shifty Backend is Ready!",
		})
	})

	// -----------------------------SERVER GRAPHQL-----------------------------
	srv := gqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: gqlResolver}))
	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))
	app.All("/query", middleware.Protected(tokenMaster), func(ctx *fiber.Ctx) error {
		userId := ctx.Locals("user_id")

		if userId != nil {
			userCtx := context.WithValue(ctx.UserContext(), "user_id", userId.(string))
			ctx.SetUserContext(userCtx)
		}
		return adaptor.HTTPHandler(srv)(ctx)
	})
	// -----------------------------AUTH GROUP----------------------------------
	auth := api.Group("/auth")
	auth.Post("/register", h.AuthHandler.RegisterLocal)
	auth.Post("/login", h.AuthHandler.LoginLocal)
	auth.Post("/send-otp", h.AuthHandler.SendOTP)

	protected := api.Group("/", middleware.Protected(tokenMaster))
	protected.Get("/pofile", h.UserHandler.UpdateAvatar)
}
