package middlewares

import (
	"strings"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func Permissions(c *fiber.Ctx) error {
	path := string(c.Request().URI().Path())
	route := strings.Split(path, "/")

	method := c.Method()

	if method == "GET" && (strings.Contains(path, "monitor") || strings.Contains(path, "swagger") || strings.Contains(path, "login") || strings.Contains(path, "register") || strings.Contains(path, "user")) {
		return c.Next()
	}

	var perm uint64
	perm = 0

	for _, uri := range route {

		switch method {
		case "GET":
			perm |= domain.Permissions["read_"+uri]

		case
			"POST",
			"PUT":
			perm |= domain.Permissions["write_"+uri]

		case "DELETE":
			perm |= domain.Permissions["delete_"+uri]
		case "PATCH":
			perm |= domain.Permissions["edit_"+uri]
		default:
		}
	}

	perm |= domain.Permissions["administrator"]

	return CheckPerms(c, perm)

}
