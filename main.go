package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/NaufalFarros/miniproject_alterra_golang/Routes"
)

func main() {
	app := fiber.New()

	Routes.routes(app)

	app.Listen(":3000")
}
