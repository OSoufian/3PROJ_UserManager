package controllers

import (
	"strconv"
	"webauthn_api/internal/domain"

	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RoleBootstrap(app fiber.Router) {

	app.Get("/:roleId", getRole)

	app.Get("/:channId", getRoles)

	app.Post("/:roleId", usersRole)

	app.Put("/:channId", createRole)

	app.Patch("/:roleId", patchRole)

	app.Delete("/:roleId", deleteRole)

}

// Get Role
// @Summary get Role
// @Description Get a specific role by an Id
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/:roleId [get]
func getRole(c *fiber.Ctx) error {
	roleId, err := strconv.ParseInt(c.Params("roleId"), 10, len(c.Params("roleId")))

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}

	role := domain.Role{}
	role.Id = uint(roleId)
	return c.Status(200).JSON(role.Get())
}

// Get Roles
// @Summary get Roles
// @Description Get all roles from a channel Id
// @Tags Roles
// @Success 200 {Role} domain.Role []
// @Failure 404
// @Router /roles/ [get]
func getRoles(c *fiber.Ctx) error {
	channel, err := utils.GetChannel(c)
	if err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}

	return c.Status(200).JSON(channel.GetRoles())

}

// Add or Remove Users Role
// @Summary add or remove user
// @Description add or remove a bulk of user roles
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/:roleId [post]
func usersRole(c *fiber.Ctx) error {
	roleId, err := strconv.ParseInt(c.Params("roleId"), 10, len(c.Params("roleId")))

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}

	role := domain.Role{}
	role.Id = uint(roleId)
	role = *role.Get()
	usersName := utils.UserRoles{}
	if err := usersName.Unmarshal(c.Body()); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	for _, v := range usersName.Usernames {
		tmpUser := domain.UserModel{}
		tmpUser.Username = v
		role.User = append(role.User, *tmpUser.Get())
	}

	role.Update()

	return c.Status(fiber.StatusAccepted).JSON(role)

}

// Create a role
// @Summary Create a role
// @Description create a role for a specify channel
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/ [put]
func createRole(c *fiber.Ctx) error {

	return c.Status(201).JSON(utils.GetRolesBody(c).Create())

}

// Patch a role
// @Summary patch a role
// @Description edit a role
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/ [patch]
func patchRole(c *fiber.Ctx) error {

	return c.Status(200).JSON(utils.GetRolesBody(c).Update().Update())

}

// delte a role
// @Summary delete a role
// @Description delete a role
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/ [delete]
func deleteRole(c *fiber.Ctx) error {
	roleId, err := strconv.ParseInt(c.Params("roleId"), 10, len(c.Params("roleId")))

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}

	role := domain.Role{}
	role.Id = uint(roleId)
	role = *role.Get()

	role.Delete()

	return c.SendStatus(fiber.StatusAccepted)

}
