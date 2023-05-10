package controllers

import (
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

type Category struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name" validate:"required"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func CreateCategory(c *fiber.Ctx) error {

	var Category Category

	if err := c.BodyParser(&Category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	errors := helper.ValidationStruct(c, Category)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	result := database.Database.Db.Create(&Category)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success",
		"data":    Category,
	})
}

func GetCategories(c *fiber.Ctx) error {
	var Categories []Category

	result := database.Database.Db.Find(&Categories)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	check := database.Database.Db.Find(&Categories).RowsAffected
	if check == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Data Category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    Categories,
	})
}

func GetCategory(c *fiber.Ctx) error {
	id := c.Query("id")
	var Category Category

	result := database.Database.Db.Where("id = ?", id).Find(&Category)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category Not Found",
		})
	}

	check := database.Database.Db.Where("id = ?", id).Find(&Category).RowsAffected
	if check == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Data Category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    Category,
	})

}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var Cat models.Category

	if err := c.BodyParser(&Cat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	result := database.Database.Db.Where("id = ?", id).Model(&Cat).Updates(&Cat)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// set response
	var Res = Category{}
	Res.ID = Cat.ID
	Res.Name = Cat.Name
	Res.Created_at = Cat.Created_at
	Res.Updated_at = Cat.Updated_at

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Update Category",
		"data":    Res,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	var Category = models.Category{}

	result := database.Database.Db.Where("id = ?", c.Params("id")).First(&Category)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category Not Found",
		})
	}

	result = database.Database.Db.Delete(&Category)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Delete Category",
	})

}
