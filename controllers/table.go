package controllers

import (
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

func CreateTable(c *fiber.Ctx) error {
	var table = models.Table{}

	if err := c.BodyParser(&table); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	err := helper.ValidationStruct(c, table)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	result := database.Database.Db.Create(&table)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success",
		"data":    table,
	})

}

func GetTables(c *fiber.Ctx) error {
	var tables []models.Table
	result := database.Database.Db.Find(&tables)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    tables,
	})

}

func UpdateTable(c *fiber.Ctx) error {

	var table models.Table
	if err := c.BodyParser(&table); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	err := helper.ValidationStruct(c, table)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	result := database.Database.Db.Model(&table).Where("id = ?", table.ID).Updates(&table)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    table,
	})

}
