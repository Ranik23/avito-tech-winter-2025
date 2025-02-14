package tokenutil

import (
	"time"

	"github.com/golang-jwt/jwt"
)


var JWTSecret = []byte("your_secret_key")

func GenerateJWT(userName string) (string, error) {
	claims := jwt.MapClaims{
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}