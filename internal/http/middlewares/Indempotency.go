package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
)

func Idempotency() fiber.Handler {
	return idempotency.New(idempotency.Config{
		Lifetime: 42 * time.Minute,
	})
}
