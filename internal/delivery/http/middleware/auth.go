package middleware

import (
	"shifty-backend/pkg/token"
	"shifty-backend/pkg/xerror"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Protected(tokenMaster *token.TokenMaster) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return xerror.Unauthorized("Missing Authorization Header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return xerror.Unauthorized("Invalid Authorization Header Format")
		}

		tokenString := parts[1]

		claims, error := tokenMaster.VerifyAccessToken(tokenString)
		if error != nil {
			return xerror.Unauthorized("Invalid or Expired Access Token")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AllowRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")

		if userRole != nil {
			return xerror.Unauthorized("Role not found in context")
		}

		roleStr, ok := userRole.(string)
		if !ok {
			return xerror.Unauthorized("Inalid role type")
		}

		for _, role := range allowedRoles {
			if role == roleStr {
				c.Next()
			}
		}
		return xerror.Forbidden("You do not have permit to access this resource")
	}
}
