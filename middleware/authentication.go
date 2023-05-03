package middleware

import (
	"strings"

	"github.com/NaufalFarros/miniproject_alterra_golang/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(c *fiber.Ctx) error {

	// get token from header
	tokenString := c.Get("Authorization")

	// dectkn, err := config.DecryptToken(tokenString)

	// if err != nil {
	// 	log.Println(err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"message": "Internal Server Error dectkn",
	// 	})
	// }

	// tokenString = dectkn

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// remove "Bearer " prefix
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	// parse token
	token, err := jwt.ParseWithClaims(tokenString, &config.JTWClaim{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	// check token validity
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// get username from token and add to context
	claims, ok := token.Claims.(*config.JTWClaim)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	c.Locals("username", claims.Username)
	return c.Next()
}
