package utils

import (
	"encoding/json"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PartialUser struct {
	Icon     string `json:"icon"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *PartialUser) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &p)
}

func CheckUserName(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")

	return c.Status(200).JSON(fiber.Map{
		"user": user.Get() != nil,
	})

}
