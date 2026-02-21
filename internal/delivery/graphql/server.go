package graphql

import (
	"context"
	"shifty-backend/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func NewGraphQLHandler(resolver *graph.Resolver) (fiber.Handler, fiber.Handler) {

	// Initialize GraphQL Server Core
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Convert between fiber and graphQL
	playgroundHandler := adaptor.HTTPHandler(playground.Handler("GraphQL Playground", "/query"))

	// Create a Handler to process the main Query
	// Wrap it to inject the UserID from the Fiber into the GraphQL Context
	queryHandler := func(c *fiber.Ctx) error {

		// Get userID from from fiber locals
		userID := c.Locals("user_id")

		if userID != nil {
			// Create a standard Go context from the Fiber request context
			ctx := context.WithValue(c.UserContext(), "user_id", userID)
			// Assign it back to the request so that the Resolver can call ctx.Value("user_id") later
			c.SetUserContext(ctx)
		}

		// Forward the request with the context to gqlgen for processing
		return adaptor.HTTPHandler(srv)(c)
	}

	return playgroundHandler, queryHandler
}
