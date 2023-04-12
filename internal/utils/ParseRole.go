package utils

import (
	"encoding/json"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PartialRole struct {
	Id          int
	Permission  int64
	Name        string
	ChannelId   int
	Description string
	Weight      int
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

	role := domain.Role{}
	role.Id = uint(partialRole.Id)
	_, err := role.Get()
	if err != nil {

	}

	role.Description = partialRole.Description
	role.Name = partialRole.Name
	role.Weight = partialRole.Weight
	role.Permission = partialRole.Permission

	channel := domain.Channel{
		Id: uint(partialRole.ChannelId),
	}

	channel.Get()
	role.Channel = channel
	role.ChannelId = int(channel.Id)

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

	return &role
}
