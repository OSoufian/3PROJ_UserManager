package middlewares

import (
	"strings"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func Permissions(c *fiber.Ctx) error {
	path := string(c.Request().URI().Path())

	route := strings.Split(path, "/")

	if strings.Contains(path, "user") {
		return c.Next()
	}

	method := c.Method()
	if method == "POST" && (strings.Contains(path, "login") || strings.Contains(path, "register")) {
		return c.Next()
	}

	if method == "GET" && (strings.Contains(path, "monitor") || strings.Contains(path, "swagger")) || strings.Contains(path, "files") || strings.Contains(path, "perms") {
		return c.Next()
	}

	var perm int64
	perm = 0

	for _, uri := range route {

		switch method {
		case "GET":
			perm |= domain.Permissions["read_"+uri]

		case "POST":
			perm |= domain.Permissions["write_"+uri]

		case "PUT":
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
