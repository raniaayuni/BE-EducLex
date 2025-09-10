package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("bG9uZ1NlY3JldEtleUFzd2RmZzEyMzQ1Njc4OQ==") // nanti diganti dengan value dari .env

func GenerateJWT(email, role string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // expired 1 hari
	})

	return token.SignedString(secret)
}
