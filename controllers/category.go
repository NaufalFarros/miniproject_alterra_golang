package controllers

import (
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {

	var Category = models.Category{}

	if err := c.BodyParser(&Category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	// cek validasi
	errors := helper.ValidationStruct(c, Category)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	Category.CreatedAt = time.Now()
	Category.UpdatedAt = time.Now()

	// save to database

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
	var Categories []models.Category

	result := database.Database.Db.Find(&Categories)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    Categories,
	})
}

func GetCategory(c *fiber.Ctx) error {
	var Category = models.Category{}

	result := database.Database.Db.Where("id = ?", c.Params("id")).First(&Category)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    Category,
	})

}

func UpdateCategory(c *fiber.Ctx) error {
	var Category = models.Category{}

	result := database.Database.Db.Where("id = ?", c.Params("id")).First(&Category)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category Not Found",
		})
	}

	if err := c.BodyParser(&Category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	// cek validasi
	errors := helper.ValidationStruct(c, Category)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	Category.UpdatedAt = time.Now()

	// save to database

	result = database.Database.Db.Save(&Category)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Update Category",
		"data":    Category,
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
