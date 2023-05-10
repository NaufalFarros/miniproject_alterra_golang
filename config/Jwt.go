package config

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var JWT_KEY = []byte(os.Getenv("JWT_KEY"))

type JWTClaim struct {
	Email  string `json:"email"`
	Roles  string `json:"roles"`
	UserID int    `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(email string, roles string, userID int) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	claims := &JWTClaim{
		Email:  email,
		Roles:  roles,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "alterra",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_KEY)
}
