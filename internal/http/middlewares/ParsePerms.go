package middlewares

import (
	"log"
	"strconv"
	"strings"
	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"
	// "fmt"

	"github.com/gofiber/fiber/v2"
)

func CheckPerms(c *fiber.Ctx, bins int64) error {
	session := utils.CheckAuthn(c)
	path := string(c.Request().URI().Path())
	route := strings.Split(path, "/")
	log.Println("session", session)
	log.Println("path", path)
	log.Println("route", route)

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
		perm |= 1118481
	} else {
		log.Println("Session not null")

		user := domain.UserModel{}
		channel := domain.Channel{}

		user.Username = session.DisplayName
		user.Get()

		log.Println("user", user)

		channPerm = 0

		if len(route) > 1 && strings.Contains(path, "video") && !strings.Contains(path, "chann") {
			log.Println("Watching")
			videoId, err := strconv.ParseInt(route[len(route)-1], 10, 64)
			if err != nil {
				log.Println("Error parsing video ID:", err)
			} else {
				log.Println("video Id", videoId)
				channel := *channel.GetByVideoId(uint(videoId))

				log.Println("channel", channel)
				log.Println("user 2", user)

				if channel.OwnerId == user.Id {
					return c.Next()
				}

				perms, err := channel.GetUserRole(user)
				if err == nil {
					log.Println("perms", perms)
					
					for _, p := range perms {
						log.Println("rolePerm", rolePerm)
						rolePerm = ^rolePerm | rolePerm&p.Permission
					}

					channPerm &= rolePerm
					log.Println("channPerm (role)", channPerm)
				} else {
					log.Println(err)
				}
			}
		}

		tempChannId, _ := strconv.ParseInt(route[len(route)-1], 10, 64)

		if (!strings.Contains(path, "undefined") && tempChannId != 0) && (strings.Contains(path, "chann") || c.Query("channId") != "") {
			log.Println("In channel")

			var (
				channId int64
				err     error
			)

			userChannel, err := user.GetChannel()

			if strings.Contains(path, "chann") && len(route) > 1 {
				log.Println("In channel with chann")
				channId, err = strconv.ParseInt(route[len(route)-1], 10, 64)
				// channId = int64(userChannel.Id)
			} else if c.Query("channId") != "" {
				log.Println("In channel with query channId")
				channId, err = strconv.ParseInt(c.Query("channId"), 10, 64)
			}

			log.Println("user Channel", userChannel)
			log.Println(uint(channId))
			log.Println("error", err)
			
			if err == nil && uint(channId) == userChannel.Id {
				return c.Next()
			}
			
			log.Println("error 2", err)
			if err != nil {
				log.Println("Error parsing channel ID:", err)
			} else {
				log.Println("Should be here")
				channel.Id = uint(channId)
				channel = *channel.Get()
				log.Println("channel 2", channel)

				channelPerms, err := channel.GetUserRole(user)
				log.Println("channPerms", channelPerms)

				if err == nil {
					// user.Permission |= perms.Permission // permissions order ASC
					for _, p := range channelPerms {
						log.Println("channPerm", channPerm)
						channPerm |= p.Permission
					}
				} else {
					log.Println(err)
				}

				log.Println("domain perm admin", domain.Permissions["admin"])
				if channPerm == domain.Permissions["admin"] {
					return c.Next()
				}
			}

			log.Println("bins", bins)
			log.Println("channPerm 2", channPerm)
			log.Println("perm 0", perm)

			bins |= channPerm
			perm |= bins

			log.Println("bins dup", bins)
			log.Println("channPerm 2 dup", channPerm)
			log.Println("perm 0 dup", perm)
		} else {
			log.Println("bins 2", bins)
			log.Println("user perms (video)", user.Permission)
			perm |= (user.Permission & bins)
			log.Println("perm 3", perm)
		}
	}

	log.Println("perm 2", perm)

	if perm != 0 {
		return c.Next()
	} else {
		if session == nil {
			log.Println("session does not exist")
		} else {
				log.Println(session.DisplayName, "is unauthorized to make >", c.Method(), "on this route >", string(c.Request().URI().Path()))
		}
		return c.Status(fiber.StatusUnauthorized).SendString("Not enought Permissions")
	}
}
