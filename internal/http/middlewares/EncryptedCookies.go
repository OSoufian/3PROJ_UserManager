package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/utils"
)

func EncryptCookie() fiber.Handler {
	return encryptcookie.New(encryptcookie.Config{
		Key: utils.UUIDv4(),
	})
}
