package main

import (
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.AdminRoutes(app)
	routes.UsersRoutes(app)
	database.DBConnect()

	app.Listen(":3000")
}
