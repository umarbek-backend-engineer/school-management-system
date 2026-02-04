package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Sign_Token(userID int, username, userRole string) (string, error) {
	jwtSecreteString := os.Getenv("JWT_SECRETE_STRING")
	jwtTimeExpires := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":       userID,
		"username":  username,
		"user_role": userRole,
	}

	if jwtTimeExpires != "" {
		duration, err := time.ParseDuration(jwtTimeExpires)
		if err != nil {
			return "", err
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		log.Println("Setting time on our own")
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Hour * 4380))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed_token, err := token.SignedString([]byte(jwtSecreteString))
	if err != nil {
		return "", err
	}

	return signed_token, nil
}
