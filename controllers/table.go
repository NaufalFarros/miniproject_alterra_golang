package controllers

import (
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

type ResponseTable struct {
	ID        uint      `json:"id"`
	Table_num string    `json:"table_number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateTable(c *fiber.Ctx) error {
	var table = models.Table{}

	if err := c.BodyParser(&table); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	// Periksa apakah nomor meja sudah ada dalam tabel
	var existingTable models.Table
	if err := database.Database.Db.Where("table_number = ?", table.Table_number).First(&existingTable).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Table number already exists",
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

	response := ResponseTable{
		ID:        table.ID,
		Table_num: table.Table_number,
		CreatedAt: table.Created_at,
		UpdatedAt: table.Updated_at,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success",
		"data":    response,
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

	var response []ResponseTable
	for _, t := range tables {
		response = append(response, ResponseTable{
			ID:        t.ID,
			Table_num: t.Table_number,
			CreatedAt: t.Created_at,
			UpdatedAt: t.Updated_at,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    response,
	})

}

func UpdateTable(c *fiber.Ctx) error {
	var table models.Table
	if err := c.BodyParser(&table); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	id := c.Params("id")
	err := helper.ValidationStruct(c, table)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	result := database.Database.Db.Model(&table).Where("id = ?", id).Updates(&table)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	// ambil data table yang sudah diupdate
	var updatedTable models.Table
	if err := database.Database.Db.Where("id = ?", id).First(&updatedTable).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	response := ResponseTable{
		ID:        updatedTable.ID,
		Table_num: updatedTable.Table_number,
		CreatedAt: updatedTable.Created_at,
		UpdatedAt: updatedTable.Updated_at,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    response,
	})
}
