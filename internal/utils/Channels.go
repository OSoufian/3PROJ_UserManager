package utils

import (
	"encoding/json"
	"strconv"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PartialChannel struct {
	Description string
	SocialLink  string
	Banner      string
	Icon        string
}

func (p *PartialChannel) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &p)
}

func ParseChannel(c *fiber.Ctx) *domain.Channel {
	session := CheckAuthn(c)

	user := domain.UserModel{}
	user.Username = session.DisplayName
	user.Get()

	partial := PartialChannel{}
	partial.Unmarshal(c.Body())

	channel := domain.Channel{}

	if partial.Banner != "" {
		channel.Banner = partial.Banner
	}
	if partial.Description != "" {
		channel.Description = partial.Description
	}

	if partial.Icon != "" {
		channel.Icon = partial.Icon
	}

	if partial.SocialLink != "" {
		channel.Icon = partial.SocialLink
	}

	return &channel
}

func GetChannel(c *fiber.Ctx) (domain.Channel, error) {
	chanId, err := strconv.ParseInt(c.Params("channId"), 10, 64)
	if err != nil {
		return domain.Channel{}, err
	}

	channel := domain.Channel{}

	channel.Id = uint(chanId)

	channel.Get()
	return channel, nil

}
