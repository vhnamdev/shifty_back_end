package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func NewRateLimiter(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "fail",
				"message": "You send request to fast! Please try again!",
			})
		},
	})
}
