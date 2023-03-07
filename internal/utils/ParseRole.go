package utils

import (
	"encoding/json"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PartialRole struct {
	Permission  uint64 `json:"permission"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserRoles struct {
	Usernames []string
}

func (ur *UserRoles) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &ur)
}

func (p *PartialRole) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &p)
}

func GetRolesBody(c *fiber.Ctx) *domain.Role {

	role := domain.Role{}

	userSession := CheckAuthn(c)
	if userSession == nil {
		return &role
	}

	channel, err := GetChannel(c)
	if err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}
	role.Channel = *channel.Get()

	return &role
}
