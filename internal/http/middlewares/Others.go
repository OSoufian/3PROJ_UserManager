package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func OthersApi(router fiber.Router) {
	//os.Getenv("ChatsAPI")
	router.Use(proxy.BalancerForward(strings.Split(os.Getenv("Others"), ";")))
}
