package controllers

import (
	"strconv"
	"strings"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func UserBootstrap(app fiber.Router) {

	app.Get("/", about)

	app.Get("/channel", getUserChannel)

	app.Get("/channel/:username", getChannelByUser)

	app.Get("/logout", logout)

	app.Post("/subscribe/:channId", nerverForget)

	app.Post("/role/:roleId", editRole)

	app.Patch("/", editUser)

	app.Delete("/", deleteUser)

	app.Delete("/cred", deleteCred)

}

// Get User
// @Summary Get about me
// @Description get all information about me
// @Tags Users
// @Success 200 {UserModel} domain.UserModel
// @Failure 404
// @Router /user [get]
func about(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user.Username = userSession.DisplayName
	return c.Status(200).JSON(user.Get())
}

// Get Channel by username
// @Summary Get channel of the user by username
// @Description get all videos of the user by username
// @Tags Channels
// @Success 200 {Channel} domain.Channel
// @Failure 404
// @Router /channel/:username [get]
func getChannelByUser(c*fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")
	
	return c.Status(200).JSON(user.GetChannel())
}

// Get Channel
// @Summary Get channel of the user
// @Description get all videos of the user
// @Tags Channels
// @Success 200 {Channel} domain.Channel
// @Failure 404
// @Router /channel [get]
func getUserChannel(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user.Username = userSession.DisplayName
	return c.Status(200).JSON(user.GetChannel())
}

// Logout
// @Summary Just Logout
// @Tags Users
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user/logout [get]
func logout(c *fiber.Ctx) error {
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.Status(200).JSON(fiber.Map{
			"message": "logout",
		})
	}
	delete(utils.Sessions, userSession.DisplayName)
	return c.Status(200).JSON(fiber.Map{
		"message": "logout",
	})
}

// Subscribe
// @Summary Subscribe
// @Tags Users
// @Description Subscribe to a channel
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user/subscribe/:channId [post]
func nerverForget(c *fiber.Ctx) error {
	userSession := utils.CheckAuthn(c)
	userIn := new(domain.UserModel)
	userIn.Username = userSession.DisplayName
	userIn = userIn.Get()

	channId, err := strconv.ParseInt(c.Params("channId"), 10, len(c.Params("channId")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	channel := domain.Channel{}

	channel.Id = uint(channId)
	channel = *channel.Get()

	if index := utils.ContainsChannel(*userIn, channel); index != -1 {
		userIn.Subscribtion = append(userIn.Subscribtion[:index], userIn.Subscribtion[index+1:]...)
	} else {
		userIn.Subscribtion = append(userIn.Subscribtion, channel)
	}

	userIn.Update()

	return c.Status(fiber.StatusAccepted).JSON(userIn)

}

// Roles
// @Summary  roles
// @Tags Users
// @Description add or remove roles to user
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user/role/:roleId [post]
func editRole(c *fiber.Ctx) error {
	userSession := utils.CheckAuthn(c)
	userIn := new(domain.UserModel)
	userIn.Username = userSession.DisplayName
	userIn = userIn.Get()

	roleId, err := strconv.ParseInt(c.Params("channId"), 10, len(c.Params("roleId")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	role := domain.Role{}

	role.Id = uint(roleId)
	role = *role.Get()

	if index := utils.HasRole(*userIn, role); index != -1 {
		userIn.Role = append(userIn.Role[:index], userIn.Role[index+1:]...)
	} else {
		userIn.Role = append(userIn.Role, role)
	}

	userIn.Update()

	return c.Status(fiber.StatusAccepted).JSON(userIn)

}

// Edit me
// @Summary  edit user
// @Tags Users
// @Description edit user information
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user [patch]
func editUser(c *fiber.Ctx) error {

	userSession := utils.CheckAuthn(c)

	userIn := new(domain.UserModel)
	userIn.Username = userSession.DisplayName
	userIn = userIn.Get()

	user := new(utils.PartialUser)
	err := user.Unmarshal(c.Body())
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	userIn.Email = user.Email
	userIn.Password = user.Password

	userIn.Update()

	return c.Status(200).JSON(user)

}

// Delete me
// @Summary  delete account
// @Tags Users
// @Description delete user account
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user [delete]
func deleteUser(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	user.Username = userSession.DisplayName

	user.Delete()
	delete(utils.Sessions, user.Username)

	return c.JSON(fiber.Map{
		"message": "deleted",
	})
}

// Delete credential
// @Summary  delete credential
// @Tags Users
// @Description delete webauthn credential
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user/cred [delete]
func deleteCred(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	user.Username = userSession.DisplayName
	user = user.Get()
	if user == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user.Incredentials = strings.Split(user.Incredentials, ";")[0]
	user.Update()

	return c.Status(200).JSON(user)
}
