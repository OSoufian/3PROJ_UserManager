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

	app.Get("/online", getOnlineUsers)

	// app.Get("/online/:videoId", getOnlineUsers)

	app.Get("/chat/:UserId", getUserById)

	app.Get("/admin/all", getAllUsers)

	app.Get("/admin/:ChannId/videos", getAllUserVideos)

	app.Patch("/admin/edit", adminEditUser)

	// app.Get("/admin/channel", getAdminUserChannel)

	app.Get("/channel", getUserChannel)

	app.Get("/admin/:username/channel", getChannelByUser)

	app.Get("/logout", logout)

	app.Post("/subscribe/:channId", nerverForget)

	app.Post("/role/:roleId", editRole)

	app.Patch("/", editUser)

	app.Delete("/", deleteUser)

	app.Delete("/:userId", deleteUserById)

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

// Get User
// @Summary Get about me
// @Description get all information about me
// @Tags Users
// @Success 200 {UserModel} domain.UserModel
// @Failure 404
// @Router /user [get]
func getUserById(c *fiber.Ctx) error {

	id := c.Params("UserId")
	UserId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err)
	}
	user := new(domain.UserModel)

	user.Id = uint(UserId)
	return c.Status(fiber.StatusAccepted).JSON(user.GetById())
}

// Get Channel by username
// @Summary Get channel of the user by username
// @Description get all video of the user by username
// @Tags Channels
// @Success 200 {Channel} domain.Channel
// @Failure 404
// @Router /user/channel/:username [get]
func getChannelByUser(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")
	user.Get()
	channel := new(domain.Channel)
	channel.OwnerId = user.Id
	channel.GetByOwner()
	return c.Status(200).JSON(channel)
}

// Get All users
// @Summary Users
// @Description get all user
// @Tags Users
// @Success 200 {Users} List Users
// @Failure 404
// @Router /all [get]
func getAllUsers(c *fiber.Ctx) error {

	if c.Query("q") != "" {
		return getUsersByCondition(c)
	}

	UsersModels := domain.UserModel{}
	users, err := UsersModels.GetAll()
	if err != nil {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	return c.Status(200).JSON(users)
}

// Get Searched userd
// @Summary users
// @Description get searched userd
// @Tags users
// @Success 200 {UserModel} List Users
// @Failure 404
// @Router /all?q="" [get]
func getUsersByCondition(c *fiber.Ctx) error {
	searchQuery := c.Query("q") // Get the search query from the query parameter "q"
	userModels := domain.UserModel{}
	users, err := userModels.GetAllCondition(searchQuery)
	if err != nil {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// Get Channel
// @Summary Get channel of the user
// @Description get all video of the user
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
	user.Get()
	channel := new(domain.Channel)
	channel.OwnerId = user.Id
	channel.GetByOwner()
	return c.Status(200).JSON(channel)
}

// // Get Channel
// // @Summary Get channel of the user
// // @Description get all video of the user
// // @Tags Channels
// // @Success 200 {Channel} domain.Channel
// // @Failure 404
// // @Router /admin/channel [get]
// func getAdminUserChannel(c *fiber.Ctx) error {
// 	user := new(domain.UserModel)
// 	userSession := utils.CheckAuthn(c)
// 	if userSession == nil {
// 		return c.SendStatus(fiber.StatusUnauthorized)
// 	}
// 	user.Username = userSession.DisplayName
// 	user.Get()
// 	channel := new(domain.Channel)
// 	channel.OwnerId = user.Id
// 	channel.GetByOwner()
// 	return c.Status(200).JSON(channel)
// }

// Logout
// @Summary Just Logout
// @Tags Users
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user/logout [get]
func logout(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.Status(200).JSON(fiber.Map{
			"message": "logout",
		})
	}

	user.Username = userSession.DisplayName
	user.Get()
	user.Online = false

	delete(utils.Sessions, userSession.DisplayName)

	user.Update()

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

	channel := domain.Channel{
		Owner:   *userIn,
		OwnerId: userIn.Id,
	}

	_, err := channel.GetByOwner()

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	roleId, err := strconv.ParseInt(c.Params("roleId"), 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	role := domain.Role{}

	role.Id = uint(roleId)
	r, err := role.Get()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	role = *r
	if index := utils.HasRole(*userIn, role); index != -1 {
		userIn.Roles = append(userIn.Roles[:index], userIn.Roles[index+1:]...)
	} else {
		userIn.Roles = append(userIn.Roles, role)
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
	session := utils.CheckAuthn(c)
	
	user := domain.UserModel{}
	user.Username = session.DisplayName
	user.Get()

	partial := new(utils.PartialUser)
	if err := partial.Unmarshal(c.Body()); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	user.Username = partial.Username

	if partial.Icon != "" {
		user.Icon = partial.Icon
	}
	user.Email = partial.Email
	user.Password = partial.Password
	user.Disable = partial.Disable

	user.Update()

	return c.Status(200).JSON(user)

}

// Edit user
// @Summary  edit user
// @Tags Users
// @Description edit user information
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user [patch]
func adminEditUser(c *fiber.Ctx) error {
	session := utils.CheckAuthn(c)
	
	adminUser := domain.UserModel{}
	adminUser.Username = session.DisplayName
	adminUser.Get()
	
	user := domain.UserModel{}
	partial := new(utils.PartialUser)
	if err := partial.Unmarshal(c.Body()); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	user.Username = partial.Username
	user.Get()

	if partial.Icon != "" {
		user.Icon = partial.Icon
	}
	user.Email = partial.Email
	user.Password = partial.Password
	user.Disable = partial.Disable

	user.Update()

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
	user.Get()

	user.Delete()
	delete(utils.Sessions, user.Username)

	return c.JSON(fiber.Map{
		"message": "deleted",
	})
}

// Delete me
// @Summary  delete account
// @Tags Users
// @Description delete user account
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user [delete]
func deleteUserById(c *fiber.Ctx) error {
	userId, er := strconv.ParseInt(c.Params("userId"), 10, 64)

	if er != nil {
		return c.Status(fiber.StatusBadRequest).JSON(er.Error())
	}

	user := new(domain.UserModel)
	user.Id = uint(userId)
	user.Get()

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
	if userSession == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user.Username = userSession.DisplayName
	user = user.Get()
	if user == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user.Incredentials = strings.Split(user.Incredentials, ";")[0]
	user.Update()

	return c.Status(200).JSON(user)
}

func getOnlineUsers(c *fiber.Ctx) error {
	keys := make([]*utils.UserSessions, 0, len(utils.Sessions))

	for k := range utils.Sessions {
		keys = append(keys, utils.Sessions[k])
	}

	return c.Status(fiber.StatusOK).JSON(keys)
}

// Get User Videos
// @Summary Videos
// @Description get all video from a user
// @Tags Videos
// @Success 200 {Videos} List of Videos
// @Failure 404
// @Router /user/admin/:ChannId/videos [get]
func getAllUserVideos(c *fiber.Ctx) error {
	id := c.Params("ChannId")
	channId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(err)
	}

	video := domain.Videos{}
	video.ChannelId = uint(channId)

	orderByParams := c.Query("orderBy")
	orderBy := strings.Split(orderByParams, ",")
	orderedVideos := video.GetAllVideosFromUser(orderBy...)
	return c.Status(fiber.StatusAccepted).JSON(orderedVideos)
}