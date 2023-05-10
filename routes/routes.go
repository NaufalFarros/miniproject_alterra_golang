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
	app.Get("/images/:imageName", func(c *fiber.Ctx) error {
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

	authAdmin := app.Group("/admin")
	authAdmin.Use(middleware.AuthorizeAdmin)
	authAdmin.Get("/profile", controllers.GetUsers)

	authAdmin.Get("/tables", controllers.GetTables)
	authAdmin.Post("/table", controllers.CreateTable)
	authAdmin.Put("/table/:id", controllers.UpdateTable)

	authAdmin.Post("/category", controllers.CreateCategory)
	authAdmin.Get("/categories", controllers.GetCategories)
	authAdmin.Get("/category", controllers.GetCategory)
	authAdmin.Put("/category/:id", controllers.UpdateCategory)
	authAdmin.Delete("/category/:id", controllers.DeleteCategory)

	authAdmin.Get("/items", controllers.GetItems)
	authAdmin.Post("/item", controllers.CreateItem)
	authAdmin.Get("/item", controllers.GetItem)
	authAdmin.Put("/item/:id", controllers.UpdateItem)
	authAdmin.Delete("/item/:id", controllers.DeleteItem)

	authAdmin.Get("/orders", controllers.GetAllOrdersUsers)
	authAdmin.Put("/orders/:id", controllers.UpdateOrderStatus)

	authUsers := app.Group("/users")
	authUsers.Use(middleware.AuthorizeUser)
	authUsers.Get("/item", controllers.GetUsers)
	authUsers.Post("/booking", controllers.CreateBookings)
	authUsers.Post("/booking-items", controllers.CreateBookingsItems)
	authUsers.Post("/booking-items-min", controllers.CreateBookingsItemsMin)
	authUsers.Post("/booking-items-submit", controllers.SubmitOrders)
}
