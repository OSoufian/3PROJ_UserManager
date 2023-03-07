package middlewares

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

func CSRF() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header: X-Csrf-Token",
		CookieName:     "csfr_",
		CookieDomain:   os.Getenv("Origins"),
		CookieHTTPOnly: true,
		KeyGenerator:   utils.UUIDv4,
		Expiration:     1 * time.Hour,
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(string(c.Request().Header.Referer()), os.Getenv("Origins"))
		},
	})
}
