package utils

import (
	"encoding/json"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PartialRole struct {
	Permission  int64  `json:"permission"`
	Name        string `json:"name"`
	ChannelId   int    `json:"channel_id"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}

type UserRoles struct {
	Usernames []string `json:"usernames"`
}

func (ur *UserRoles) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &ur)
}

func (p *PartialRole) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &p)
}

func GetRolesBody(c *fiber.Ctx) *domain.Role {

	partialRole := PartialRole{}
	partialRole.Unmarshal(c.Body())

	role := domain.Role{
		Permission:  partialRole.Permission,
		ChannelId:   partialRole.ChannelId,
		Weight:      partialRole.Weight,
		Description: partialRole.Description,
		Name:        partialRole.Name,
	}

	userSession := CheckAuthn(c)
	if userSession == nil {
		return &role
	}
	user := domain.UserModel{
		Username: userSession.DisplayName,
	}
	user.Get()

	if user.Permission&domain.Permissions["administrator"] != domain.Permissions["administrator"] {
		role.Permission &= ^domain.Permissions["administrator"]
	}

	channel, err := GetChannel(c)
	if err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}
	role.Channel = *channel.Get()

	return &role
}
