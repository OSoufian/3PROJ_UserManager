package controllers

import (
	"webauthn_api/internal/domain"
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
// @Router /channel/:channId [get]
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
// @Router /channel/:channId [put]
func createChannel(c *fiber.Ctx) error {
	channel := utils.ParseChannel(c)
	session := utils.CheckAuthn(c)
	user := domain.UserModel{}
	user.Username = session.DisplayName
	channel.OwnerId = user.Get().Id
	channel, err := channel.GetByOwner()
	if user.Get() == nil || err == nil {
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
// @Router /channel/:channId [patch]
func patchChannel(c *fiber.Ctx) error {
	session := utils.CheckAuthn(c)
	
	user := domain.UserModel{}
	user.Username = session.DisplayName
	user.Get()

	partial := new(utils.PartialChannel)
	if err := partial.Unmarshal(c.Body()); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	channel := domain.Channel{}
	channel.Id = uint(partial.Id)
	channel.Get()

	if partial.Description != "" {
		channel.Description = partial.Description
	}
	if partial.Name != "" {
		channel.Name = partial.Name
	}
	if partial.Banner != "" {
		channel.Banner = partial.Banner
	}
	if partial.Icon != "" {
		channel.Icon = partial.Icon
	}
	if partial.SocialLink != "" {
		channel.SocialLink = partial.SocialLink
	}

	channel.Update()

	return c.Status(fiber.StatusAccepted).JSON(channel)
}

// Delete Channel
// @Summary get channel
// @Description get channel by id
// @Tags Channels
// @Success 200 {Channels} domain.Channel
// @Failure 404
// @Router /channel/:channId [delete]
func deleteChannel(c *fiber.Ctx) error {

	channel, err := utils.GetChannel(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	// channel.DeleteAllVideos()

	return c.Status(fiber.StatusAccepted).JSON(channel.Delete())
}
