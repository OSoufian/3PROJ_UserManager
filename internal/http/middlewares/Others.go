package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func OthersApi(router fiber.Router) {
	//os.Getenv("ChatsAPI")

	router.Use(func(c *fiber.Ctx) error {
		for _, server := range strings.Split(os.Getenv("Others"), ";") {
			// if !strings.HasPrefix(server, "http") || !strings.HasPrefix(server, "ws") { -> Marche pas.
			if !strings.HasPrefix(server, "http") {
				server = "http://" + server
			}

			c.Request().Header.Add("X-Real-IP", c.IP())
			err := proxy.Do(c, server+c.OriginalURL())
			if c.Response().StatusCode() != 404 {

				return err
			}

		}
		return c.SendStatus(404)
	})

}
