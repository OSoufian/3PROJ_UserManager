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
	// app.Use(middlewares.EncryptCookie())

	app.Use(middlewares.Permissions)

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

	controllers.PermissionBootstrap(app.Group("/perms"))

	controllers.RoleBootstrap(app.Group("/roles"))

	controllers.ChannelBootstrap(app.Group("/channel"))

	// middlewares.OthersApi(app.Group("/"))

	return app
}
