package utils

/*
*  return true only if a session contains AAGUID
 */

import (
	"strings"
	"webauthn_api/internal/domain"
	"time"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/duo-labs/webauthn/webauthn"
)


type UserSessions struct {
	SessionData *webauthn.SessionData `json:"-"`
	SessionCred *webauthn.Credential  `json:"-"`
	DisplayName string
	Jwt         string
	Expiration  time.Duration  `json:"-"`
	Online		bool   `json:"online"`
	videoId		int64  `json:"videoId"`
}

var Sessions map[string]*UserSessions

func (session UserSessions) DeleteAfter() {

	
	time.Sleep(session.Expiration * time.Second)

	log.Printf("user delete")

	user := domain.UserModel{}
	user.Username = session.DisplayName

	userModel := user.Get()

	if user.Password == "" && user.Incredentials == "" {
		userModel.Delete()
	}

	delete(Sessions, session.DisplayName)
}

func CheckAuthn(c *fiber.Ctx) *UserSessions {
	value, ok := c.GetReqHeaders()["Authorization"]
	if !ok {
		return nil
	}
	authType := strings.Split(value, " ")
	if authType[0] != "Bearer" || len(authType) < 2 {
		return nil
	}

	auth := authType[1]

	// log.Println(authType)

	for _, v := range Sessions {

		if CheckJWT(v, auth) {
			return v
		}

	}
	return nil
}
