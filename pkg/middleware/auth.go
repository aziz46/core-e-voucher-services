package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates API key from header
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing API key",
			})
		}

		// TODO: Validate API key against tenant database
		c.Locals("tenant_id", "tenant_default")
		return c.Next()
	}
}

// RateLimitMiddleware applies rate limiting (3 requests per second per partner)
func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement Redis-based rate limiting
		return c.Next()
	}
}
