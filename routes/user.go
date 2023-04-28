package routes

import (
	"github.com/gofiber/fiber/v2"
)

func UsersRoutes(app *fiber.App) {
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.SendString("Hello, User 👋!")
	})
}
