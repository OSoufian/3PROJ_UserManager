package middlewares

import (
	"log"
	"strconv"
	"strings"
	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func CheckPerms(c *fiber.Ctx, bins int64) error {
	session := utils.CheckAuthn(c)
	path := string(c.Request().URI().Path())
	route := strings.Split(path, "/")

	if len(route) == 1 && strings.Contains(path, "files") {
		return c.Next()
	}

	var (
		perm      int64
		channPerm int64
		rolePerm  int64
	)
	perm = 0

	if session == nil {
		perm |= 0

	} else {
		user := domain.UserModel{}
		channel := domain.Channel{}

		user.Username = session.DisplayName
		user.Get()

		channPerm = 0

		if len(route) > 1 && (strings.Contains(path, "watch")) {
			videoId, err := strconv.ParseInt(route[len(route)-1], 10, 64)
			if err != nil {
				log.Println("Error parsing video ID:", err)
			} else {
				channel := *channel.GetByVideoId(uint(videoId))

				if channel.OwnerId == user.Id {
					return c.Next()
				}

				perms, err := channel.GetUserRole(user)
				if err == nil {
					// user.Permission |= perms.Permission // permissions order ASC

					for _, p := range perms {
						rolePerm = ^rolePerm | rolePerm&p.Permission
					}

					channPerm |= rolePerm
				} else {
					log.Println(err)
				}
			}
		}

		if strings.Contains(path, "channel") || c.Params("channId") != "" || c.Query("channId") != "" {
			var (
				channId int64
				err     error
			)

			if strings.Contains(path, "channel") && len(route) > 1 {
				channId, err = strconv.ParseInt(route[len(route)-1], 10, 64)
			} else if c.Query("channId") != "" {
				channId, err = strconv.ParseInt(c.Query("channId"), 10, 64)
			} else {
				channId, err = strconv.ParseInt(c.Params("channId"), 10, 64)
			}

			if err != nil {
				log.Println("Error parsing channel ID:", err)
			} else {
				channel.Id = uint(channId)
				channel = *channel.Get()

				if channel.OwnerId == user.Id {
					return c.Next()
				}

				channelPerms, err := channel.GetUserRole(user)

				if err == nil {
					// user.Permission |= perms.Permission // permissions order ASC
					for _, p := range channelPerms {
						channPerm |= p.Permission
					}
				} else {
					log.Println(err)
				}
			}
		}

		bins |= channPerm
		perm |= (user.Permission & bins)
	}

	if perm != 0 {
		return c.Next()
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
}
