package utils

import (
	"log"

	"math/rand"

	"webauthn_api/internal/domain"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
)

var sampleSecretKey = []byte(generateKey(20))

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateKey(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CreateJWT(session domain.UserSessions) (string, error) {
	var authToken string
	if session.SessionData == nil {
		authToken = uuid.NewString()

	} else {
		authToken = string(session.SessionData.UserID)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  session.DisplayName,
		"AuthToken": authToken,
	})
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckJWT(session *domain.UserSessions, tokenString string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		log.Println(err.Error())

		return false
	}
	if token.Method != jwt.SigningMethodHS256 {
		log.Println(token.Valid)
		return false
	}

	i := 1 << 0
	for _, val := range claims {
		if val == session.DisplayName {
			i |= 1 << 1
		}
		if val == string(session.SessionData.UserID) {
			i |= 1 << 0
		}
	}

	return i == 3

}
