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

	// Get the JWT token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.New("Missing Authorization header")
	}
	tokenString := authHeader[7:]

	// remove "Bearer " prefix
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	// parse token
	token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	// check token validity
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Get the "roles" claim from the token
	claims, ok := token.Claims.(*config.JWTClaim)
	if !ok || !token.Valid {
		return errors.New("Invalid token claims")
	}
	role := claims.Roles
	if role == "" {
		return errors.New("Roles claim is missing or invalid")
	}

	// Check if the user has the "user" role
	hasUserRole := false
	if role == "user" {
		hasUserRole = true
	}

	if !hasUserRole {
		return errors.New("Unauthorized access")
	}
	// Set the user ID as a context value
	userID := claims.UserID
	c.Locals("userID", userID)

	// User is authorized, call next middleware
	return c.Next()

}

func AuthorizeAdmin(c *fiber.Ctx) error {

	// Get the JWT token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.New("Missing Authorization header")
	}
	tokenString := authHeader[7:]

	// remove "Bearer " prefix
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	// parse token
	token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	// check token validity
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Get the "roles" claim from the token
	claims, ok := token.Claims.(*config.JWTClaim)
	if !ok || !token.Valid {
		return errors.New("Invalid token claims")
	}
	role := claims.Roles
	if role == "" {
		return errors.New("Roles claim is missing or invalid")
	}

	// Check if the user has the "admin" role
	hasAdminRole := false
	if role == "admin" {
		hasAdminRole = true
	}

	if !hasAdminRole {
		return errors.New("Unauthorized access")
	}

	// User is authorized, call next middleware
	return c.Next()

}

func GetUsersFromToken(c *fiber.Ctx) (string, error) {
	// Get the JWT token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Missing Authorization header")
	}
	tokenString := authHeader[7:]

	// Remove "Bearer " prefix
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	// Check token validity
	if err != nil {
		return "", errors.New("Invalid token")
	}

	// Get the "sub" claim from the token
	claims, ok := token.Claims.(*config.JWTClaim)
	if !ok || !token.Valid {
		return "", errors.New("Invalid token claims")
	}
	userID := claims.Subject
	if userID == "" {
		return "", errors.New("User ID claim is missing or invalid")
	}

	return userID, nil
}
