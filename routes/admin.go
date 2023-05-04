package routes

import (
	"fmt"
	"os"

	"github.com/NaufalFarros/miniproject_alterra_golang/controllers"
	"github.com/NaufalFarros/miniproject_alterra_golang/middleware"
	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App) {
	app.Post("/login", controllers.Login)
	app.Post("/register", controllers.Register)
	app.Post("/logout", controllers.Logout)

	auth := app.Group("/admin")
	auth.Use(middleware.JWTAuthMiddleware)
	auth.Get("/profile", controllers.GetUsers)
	auth.Post("/category", controllers.CreateCategory)
	auth.Get("/category", controllers.GetCategories)
	auth.Get("/category/:id", controllers.GetCategory)
	auth.Put("/category/:id", controllers.UpdateCategory)
	auth.Delete("/category/:id", controllers.DeleteCategory)

	auth.Get("/items", controllers.GetItems)
	auth.Post("/item", controllers.CreateItem)
	auth.Get("/item", controllers.GetItem)

	auth.Get("/images/:imageName", func(c *fiber.Ctx) error {
		imageName := c.Params("imageName")
		imagePath := "./image/" + imageName
		fmt.Println("Ianmge Path :", imagePath)
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Image not found",
			})
		}
		return c.SendFile(imagePath)
	})

}
