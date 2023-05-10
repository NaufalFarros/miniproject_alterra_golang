package middleware

import (
	"errors"
	"os"
	"strings"

	"github.com/NaufalFarros/miniproject_alterra_golang/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() []byte {
	secretKey := os.Getenv("JWT_KEY")
	if secretKey == "" {
		panic("JWT_SECRET_KEY is not set in environment variables")
	}
	return []byte(secretKey)
}

func AuthorizeUser(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.New("missing Authorization header")
	}
	tokenString := authHeader[7:]

	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	claims, ok := token.Claims.(*config.JWTClaim)
	if !ok || !token.Valid {
		return errors.New("invalid token claims")
	}
	role := claims.Roles
	if role == "" {
		return errors.New("roles claim is missing or invalid")
	}

	hasUserRole := false
	if role == "user" {
		hasUserRole = true
	}

	if !hasUserRole {
		return errors.New("unauthorized access")
	}
	userID := claims.UserID
	c.Locals("userID", userID)

	return c.Next()

}

func AuthorizeAdmin(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.New("missing Authorization header")
	}
	tokenString := authHeader[7:]

	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims, ok := token.Claims.(*config.JWTClaim)
	if !ok || !token.Valid {
		return errors.New("invalid token claims")
	}
	role := claims.Roles
	if role == "" {
		return errors.New("roles claim is missing or invalid")
	}

	hasAdminRole := false
	if role == "admin" {
		hasAdminRole = true
	}

	if !hasAdminRole {
		return errors.New("unauthorized access")
	}

	return c.Next()

}

func GetUsersFromToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	tokenString := authHeader[7:]

	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*config.JWTClaim)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}
	userID := claims.Subject
	if userID == "" {
		return "", errors.New("user ID claim is missing or invalid")
	}

	return userID, nil
}
