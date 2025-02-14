package tokenutil

import (
	"time"

	"github.com/golang-jwt/jwt"
)


type Token interface {
	GenerateJWT(userName string) (string, error)
	ParseJWT(tokenString string) (*jwt.Token, error)
}

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

func (j *JWTService) GenerateJWT(userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(j.secret)
}

func (j *JWTService) ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.secret, nil
	})
}

var JWTSecret = []byte("your_secret_key")

func GenerateJWT(userName string) (string, error) {
	claims := jwt.MapClaims{
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}