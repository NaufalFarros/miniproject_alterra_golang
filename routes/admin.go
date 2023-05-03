package routes

import (
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

	auth.Post("/item", controllers.CreateItem)
}
