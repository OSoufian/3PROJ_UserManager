package controllers

import (
	"encoding/json"
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type partialPerm struct {
	Name string
	Bin  uint64
}

func PermissionBootstrap(app fiber.Router) {

	app.Get("/:name", getPermissions)
	app.Get("/", getAllPerms)
}

// Get Permission
// @Summary get permissions
// @Description get any permissions every where
// @Tags Perms
// @Success 200 {Permissions} partialPerm
// @Failure 404
// @Router /perms [get]
func getAllPerms(c *fiber.Ctx) error {
	j, _ := json.Marshal(domain.Permissions)
	return c.Status(200).Send(j)
}

// Get Permission
// @Summary get permissions
// @Description get any permissions every where
// @Tags Perms
// @Success 200 {Permissions} partialPerm
// @Failure 404
// @Router /:name [get]
func getPermissions(c *fiber.Ctx) error {
	name := c.Params("name")
	perm := partialPerm{}
	perm.Name = name
	perm.Bin = domain.Permissions[name]

	return c.Status(200).JSON(&fiber.Map{
		"name":  name,
		"value": perm,
	})
}
