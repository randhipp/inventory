package services

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/randhipp/inventory/models"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func (s AuthService) GetToken(user models.User) (tokenString string, err error) {
	var secret []byte
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret = []byte(os.Getenv("HMAC_SECRET"))

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":  user.Name,
		"email": user.Email,
		"id":    user.ID,
		"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(secret)
	if err != nil {
		log.Fatal("JWT SignedString ERR")
	}

	return tokenString, nil
}
