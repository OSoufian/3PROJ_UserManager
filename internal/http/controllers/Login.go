package controllers

import (
	"bytes"
	"log"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
)

func LoginBoostrap(app fiber.Router) {
	app.Post("/start/:username", loginStart)

	app.Post("/end/:username", loginEnd)

	app.Post("/password/:username", loginPassword)
}

// Begin Login
// @Summary begin Login
// @Description begin the webauthn login and update the user credential to session and database
// @Tags Logins
// @Success 200 {Options} webauthn.Options
// @Failure 404
// @Router /login/start/:username [post]
func loginStart(c *fiber.Ctx) error {

	user := new(domain.UserModel)
	user.Username = c.Params("username")
	if !user.Find() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}
	user.ParseCredentials()

	options, sessionData, err := utils.Web.BeginLogin(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	session := new(domain.UserSessions)

	session.SessionData = sessionData
	session.DisplayName = user.Username
	session.Expiration = 60 * 1000
	go session.DeleteAfter(utils.Sessions)
	utils.Sessions[user.Username] = session

	return c.JSON(options)

}

// End Login
// @Summary end Login
// @Description end the webauthn login and update the user to session and database
// @Tags Logins
// @Success 200 {JWTtoken} JWT token
// @Failure 404
// @Router /login/end/:username [post]
func loginEnd(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")
	if !user.Find() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}
	user.ParseCredentials()

	session, ok := utils.Sessions[c.Params("username")]

	if !ok {

		log.Println("session Not exist")
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "session Not exist",
		})
	}

	parsedCredential, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	creds, err := utils.Web.ValidateLogin(user, *session.SessionData, parsedCredential)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}
	session.SessionCred = creds
	session.Expiration = 24 * 3600 * 2
	token, err := utils.CreateJWT(*session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	session.Jwt = token
	go session.DeleteAfter(utils.Sessions)

	utils.Sessions[user.Username] = session
	user.Credentials = append(user.Credentials, *creds)

	user.SaveCredentials()

	return c.JSON(fiber.Map{
		"token": session.Jwt,
	})
}

// Password Login
// @Summary password Login
// @Description set a password login for users
// @Tags Logins
// @Success 200 {JWTtoken} JWT token
// @Failure 404
// @Router /login/password/:username [post]
func loginPassword(c *fiber.Ctx) error {

	user := new(domain.UserModel)
	user.Username = c.Params("username")
	if !user.Find() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}

	userBody := new(domain.UserModel)

	if err := c.BodyParser(&userBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	if userBody.Password != user.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": "Not Authorize",
		})
	}

	session := new(domain.UserSessions)

	session.DisplayName = user.Username
	token, err := utils.CreateJWT(*session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	session.Jwt = token
	session.Expiration = 24 * 3600 * 2
	go session.DeleteAfter(utils.Sessions)

	utils.Sessions[user.Username] = session

	return c.JSON(fiber.Map{
		"token": session.Jwt,
	})
}
