package middlewares

import (
	"log"
	"strconv"
	"strings"
	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func hasPermission(perm int64, wantedPerm int64) bool {
    // Decimal values to check against: 16707310
    checkValue := int64(wantedPerm)

    // Convert decimal values to binary
    permBinary := strconv.FormatInt(perm, 2)
    checkValueBinary := strconv.FormatInt(checkValue, 2)

    // Pad binary strings with leading zeros to ensure equal length
    maxLength := len(permBinary)
    if len(checkValueBinary) > maxLength {
        maxLength = len(checkValueBinary)
    }
    permBinary = fmt.Sprintf("%0*s", maxLength, permBinary)
    checkValueBinary = fmt.Sprintf("%0*s", maxLength, checkValueBinary)

    // Perform the binary check using the bitwise AND operator
    for i := 0; i < len(permBinary); i++ {
        if permBinary[i] == '1' && checkValueBinary[i] == '1' {
            return true
        }
    }

    return false
}

func CheckAdminPerms(c *fiber.Ctx, bins int64) error {
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

	if session != nil {
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

					for _, p := range perms {
						rolePerm = ^rolePerm | rolePerm&p.Permission
					}

					channPerm &= rolePerm
				} else {
					log.Println(err)
				}
			}
		}

		if strings.Contains(path, "channId") || c.Query("channId") != "" {
			var (
				channId int64
				err     error
			)

			if strings.Contains(path, "channId") && len(route) > 1 {
				channId, err = strconv.ParseInt(route[len(route)-1], 10, 64)
			} else if c.Query("channId") != "" {
				channId, err = strconv.ParseInt(c.Query("channId"), 10, 64)
			}

			userChannel, err := user.GetChannel()

			if err == nil && uint(channId) == userChannel.Id {
				return c.Next()
			}

			if err != nil {
				log.Println("Error parsing channel ID:", err)
			} else {
				channel.Id = uint(channId)
				channel = *channel.Get()

				channelPerms, err := channel.GetUserRole(user)

				if err == nil {
					// user.Permission |= perms.Permission // permissions order ASC
					for _, p := range channelPerms {
						channPerm |= p.Permission
					}
				} else {
					log.Println(err)
				}

				if channPerm == domain.Permissions["admin"] {
					return c.Next()
				}
			}

			bins |= channPerm
			perm |= bins
		} else {
			perm |= (user.Permission & bins)
		}
	}

	if perm != 0 && ( hasPermission(perm, 16707310) || hasPermission(perm, 16777216) || hasPermission(perm, 33554432)){
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
		perm |= 69905
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

					for _, p := range perms {
						rolePerm = ^rolePerm | rolePerm&p.Permission
					}

					channPerm &= rolePerm
				} else {
					log.Println(err)
				}
			}
		}

		if strings.Contains(path, "channId") || c.Query("channId") != "" {
			var (
				channId int64
				err     error
			)

			if strings.Contains(path, "channId") && len(route) > 1 {
				channId, err = strconv.ParseInt(route[len(route)-1], 10, 64)
			} else if c.Query("channId") != "" {
				channId, err = strconv.ParseInt(c.Query("channId"), 10, 64)
			}

			userChannel, err := user.GetChannel()

			if err == nil && uint(channId) == userChannel.Id {
				return c.Next()
			}

			if err != nil {
				log.Println("Error parsing channel ID:", err)
			} else {
				channel.Id = uint(channId)
				channel = *channel.Get()

				channelPerms, err := channel.GetUserRole(user)

				if err == nil {
					// user.Permission |= perms.Permission // permissions order ASC
					for _, p := range channelPerms {
						channPerm |= p.Permission
					}
				} else {
					log.Println(err)
				}

				if channPerm == domain.Permissions["admin"] {
					return c.Next()
				}
			}

			bins |= channPerm
			perm |= bins
		} else {
			perm |= (user.Permission & bins)
		}
	}

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
