package http

import (
	"webauthn_api/internal/http/controllers"

	_ "webauthn_api/docs"
	"webauthn_api/internal/http/middlewares"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func healthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}

func Http() *fiber.App {
	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
		BodyLimit:         100 * 1024 * 1024 * 1024,
	})

	app.Use(middlewares.CORS())
	// app.Use(middlewares.Idempotency())

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// app.Use(middlewares.CSRF())
	// app.Use(encryptcookie.New(encryptcookie.Config{
	// 	Key: fiberUtils.UUIDv4(),
	// }))

	// app.Get("/checkUser/:username", CheckUserName)
	app.Get("/", healthCheck)
	app.Get("/swagger/*", fiberSwagger.FiberWrapHandler())

	app.Get("/monitor", monitor.New(monitor.Config{
		Title: "Login Register Monitor",
	}))

	//app routes
	controllers.RegisterBootstrap(app.Group("/register"))

	controllers.LoginBoostrap(app.Group("/login"))

	controllers.UserBootstrap(app.Group("/user"))

	/* app.Use(func(c *fiber.Ctx) error {
		route := strings.Split(string(c.Request().URI().Path()), "/")[1]
		if route == "upload" {
			route = "video"

			if c.Method() == "PUT" {
				partial := new(utils.PartialVideo)

				partial.Unmarshal(c.Body())

				session := utils.CheckAuthn(c)

				user := new(domain.UserModel)
				user.Username = session.DisplayName
				user.Get()

				channel := new(domain.Channel)
				channel.OwnerId = user.Id

				channel.GetByOwer()

				video := new(utils.PartialCreateVideo)

				video.ChannelId = uint64(channel.Id)
				video.Description = partial.Description
				video.Name = partial.Name
				video.Icon = partial.Icon
				bittes, _ := json.Marshal(video)

				c.Request().SetBody(bittes)

			}
		}

		if c.Method() == "GET" {

			if strings.Contains(string(c.Request().URI().Path()), "monitor") || strings.Contains(string(c.Request().URI().Path()), "swagger") {
				return c.Next()
			}

			return middlewares.CheckPerms(c, domain.Permissions["read_"+route]|domain.Permissions["administrator"])

		} else if c.Method() == "PUT" {
			return middlewares.CheckPerms(c, domain.Permissions["write_"+route]|domain.Permissions["administrator"])

		} else if c.Method() == "POST" {
			return middlewares.CheckPerms(c, domain.Permissions["edit_"+route]|domain.Permissions["administrator"])

		} else if c.Method() == "PATCH" {
			return middlewares.CheckPerms(c, domain.Permissions["edit_"+route]|domain.Permissions["administrator"])
		} else if c.Method() == "DELETE" {
			return middlewares.CheckPerms(c, domain.Permissions["delete_"+route]|domain.Permissions["administrator"])

		} else {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}) */

	controllers.PermissionBootstrap(app.Group("/perms"))

	controllers.RoleBootstrap(app.Group("/roles"))

	controllers.ChannelBootstrap(app.Group("/channel"))

	middlewares.OthersApi(app.Group("/"))

	return app
}
