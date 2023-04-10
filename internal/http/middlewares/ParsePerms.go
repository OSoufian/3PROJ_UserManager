package middlewares

import (
	"log"
	"strconv"
	"strings"
	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func CheckPerms(c *fiber.Ctx, bins uint64) error {
	session := utils.CheckAuthn(c)
	path := string(c.Request().URI().Path())
	route := strings.Split(path, "/")

	var (
		perm uint64
	)

	if session == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := domain.UserModel{}
	user.Username = session.DisplayName
	user.Get()

	if strings.Contains(path, "watch") {
		videoId, err := strconv.ParseInt(route[len(route)-1], 10, 64)
		if err != nil {
			log.Println("Error parsing video ID:", err)
		} else {
			channel := domain.Channel{}
			channel = *channel.GetByVideoId(uint(videoId))
			perms, err := channel.GetUserRole(user)
			if err == nil {
				user.Permission |= perms.Permission
			} else {
				log.Println(err)
			}
		}

	}

	perm |= (user.Permission & bins)

	if perm != 0 {
		return c.Next()
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
}
