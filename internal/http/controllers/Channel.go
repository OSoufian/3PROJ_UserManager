package controllers

import (
	"github.com/gofiber/fiber/v2"
	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"
)

func ChannelBootstrap(app fiber.Router) {

	app.Get("/:channId", getChannel)
	app.Get("/:channId/video", getVideos)
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

// Get Video
// @Summary Get all video
// @Description get all video of the user
// @Tags Videos
// @Success 200 {Videos} domain.Videos
// @Failure 404
// @Router /channels/:channId/video [get]
func getVideos(c *fiber.Ctx) error {
	channel, err := utils.GetChannel(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(200).JSON(channel.GetAllVideos())

}

// Create Channel
// @Summary get channel
// @Description get channel by id
// @Tags Channels
// @Success 200 {Channels} domain.Channel
// @Failure 404
// @Router /channels/:channId [put]
func createChannel(c *fiber.Ctx) error {
	channel := utils.ParseChannel(c)
	session := utils.CheckAuthn(c)
	user := domain.UserModel{}
	user.Username = session.DisplayName
	channel.OwnerId = user.Get().Id
	if user.Get() == nil || channel.GetByOwner() != nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Status(fiber.StatusAccepted).JSON(channel.Create())
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
