package middlewares

import (
	"strconv"
	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func CheckPerms(c *fiber.Ctx, bins uint64) error {
	session := utils.CheckAuthn(c)
	var (
		perm uint64
	)

	if session == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := domain.UserModel{}
	user.Username = session.DisplayName
	user.Get()

	perm |= (user.Permission & bins)

	if c.Params("channId") != "" {

		chanId, err := strconv.ParseInt(c.Params("channId"), 10, len(c.Params("channId")))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		channel := domain.Channel{}

		channel.Id = uint(chanId)
		roles, err := channel.GetUserRole(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		if channel.OwnerId == user.Id {
			perm |= 1 << 63
		}

		for _, r := range roles {
			if r.Permission&bins != 0 {
				perm |= r.Permission & bins
				break
			}
		}

	}

	if perm != 0 {
		return c.Next()
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
}
