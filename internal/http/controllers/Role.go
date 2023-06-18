package controllers

import (
	"strconv"
	"webauthn_api/internal/domain"

	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"log"
)

func RoleBootstrap(router fiber.Router) {

	router.Get("/:roleId", getRole)

	router.Get("/channel/:channId", getRoles)

	router.Post("/:roleId", usersRole)

	router.Put("/:channId", createRole)

	router.Patch("/:roleId", patchRole)

	router.Delete("/:roleId", deleteRole)

}

// Get Role
// @Summary get Role
// @Description Get a specific role by an Id
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/:roleId [get]
func getRole(c *fiber.Ctx) error {
	roleId, err := strconv.ParseInt(c.Params("roleId"), 10, 64)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}

	role := domain.Role{}
	role.Id = uint(roleId)
	roles, _ := role.Get()
	return c.Status(200).JSON(roles)
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

	roles, err := channel.GetRoles()

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(200).JSON(roles)

}

// Add or Remove Users Role
// @Summary add or remove user
// @Description add or remove a bulk of user roles
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/:roleId [post]
func usersRole(c *fiber.Ctx) error {
	roleId, err := strconv.ParseInt(c.Params("roleId"), 10, 64)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}

	role := domain.Role{}
	role.Id = uint(roleId)
	r, _ := role.Get()
	role = *r
	usersName := utils.UserRoles{}
	if err := usersName.Unmarshal(c.Body()); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	for _, v := range usersName.Usernames {
		tmpUser := domain.UserModel{}
		tmpUser.Username = v
		role.Users = append(role.Users, *tmpUser.Get())
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
// @Router /roles/:channId [put]
func createRole(c *fiber.Ctx) error {

	role, err := utils.GetRolesBody(c).Create()

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(200).JSON(role)

}

// Patch a role
// @Summary patch a role
// @Description edit a role
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/ [patch]
func patchRole(c *fiber.Ctx) error {
	role := utils.GetRolesBody(c)

	log.Println("patch start", role)

	if _, err := role.Update(); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(role)

}

// delte a role
// @Summary delete a role
// @Description delete a role
// @Tags Roles
// @Success 200 {Role} domain.Role
// @Failure 404
// @Router /roles/ [delete]
func deleteRole(c *fiber.Ctx) error {
	roleId, err := strconv.ParseInt(c.Params("roleId"), 10, 64)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err.Error())
	}

	role := domain.Role{}
	role.Id = uint(roleId)
	r, _ := role.Get()
	role = *r

	role.Delete()

	return c.SendStatus(fiber.StatusAccepted)

}
