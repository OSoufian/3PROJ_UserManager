package utils

import (
	"encoding/json"
	"strconv"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PartialChannel struct {
	Id          int	   `json:"Id"`
	OwnerId     int    `json:"OwnerId"`
	Description string `json:"Description"`
	Name        string `json:"Name"`
	Banner      string `json:"Banner"`
	Icon        string `json:"Icon"`
	SocialLink  string `json:"SocialLink"`
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
	channel.Id = uint(partial.Id)
	channel.Get()

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
