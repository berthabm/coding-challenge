package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/interseguros/challenge/api-go/internal/utils"
)

// RequestLogger logs method, path, status, and duration for every request.
func RequestLogger(logger utils.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		logger.Info("http request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}
}
