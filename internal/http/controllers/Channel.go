package controllers

import (
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func ChannelBootstrap(app fiber.Router) {

	app.Get("/:channId", getChannel)
	app.Put("/", createChannel)
	app.Patch("/:channId", patchChannel)
	app.Delete("/:channId", deleteChannel)

}

// Get Channel
// @Summary get channel
// @Description get channel by id
// @Tags Channels
// @Success 200 {Channels} domain.Channel
// @Failure 404
// @Router /channels/:channId [get]
func getChannel(c *fiber.Ctx) error {
	channel, err := utils.GetChannel(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(200).JSON(channel)

}

// Create Channel
// @Summary get channel
// @Description get channel by id
// @Tags Channels
// @Success 200 {Channels} domain.Channel
// @Failure 404
// @Router /channels/:channId [put]
func createChannel(c *fiber.Ctx) error {

	return c.Status(fiber.StatusAccepted).JSON(utils.ParseChannel(c).Create())
}

// Patch Channel
// @Summary get channel
// @Description get channel by id
// @Tags Channels
// @Success 200 {Channels} domain.Channel
// @Failure 404
// @Router /channels/:channId [patch]
func patchChannel(c *fiber.Ctx) error {

	return c.Status(fiber.StatusAccepted).JSON(utils.ParseChannel(c).Update())
}

// Delete Channel
// @Summary get channel
// @Description get channel by id
// @Tags Channels
// @Success 200 {Channels} domain.Channel
// @Failure 404
// @Router /channels/:channId [delete]
func deleteChannel(c *fiber.Ctx) error {

	channel, err := utils.GetChannel(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusAccepted).JSON(channel.Delete())
}
