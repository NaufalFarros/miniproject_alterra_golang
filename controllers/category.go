package controllers

import (
	"strconv"
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

type Category struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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

	Category.Created_at = time.Now()
	Category.Updated_at = time.Now()

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
	var Category Category
	id := c.Query("id")
	// get one data

	result := database.Database.Db.Where("id = ?", id).Find(&Category)

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
	var Category Category

	if err := c.BodyParser(&Category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}
	// get ID from URL params and set to Category.ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	Category.ID = uint(id)
	updates := map[string]interface{}{
		"name": Category.Name,
	}

	// cek validasi
	errors := helper.ValidationStruct(c, Category)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	Category.UpdatedAt = time.Now()

	// save to database

	result := database.Database.Db.Model(&Category).Where("id = ?", id).Updates(updates)

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
