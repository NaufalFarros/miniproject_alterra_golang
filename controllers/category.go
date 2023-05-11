package controllers

import (
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

type UpdateCategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateCategory(c *fiber.Ctx) error {

	var Category models.Category

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

	response := UpdateCategoryResponse{
		ID:        Category.ID,
		Name:      Category.Name,
		CreatedAt: Category.Created_at,
		UpdatedAt: Category.Updated_at,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success",
		"data":    response,
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

	check := database.Database.Db.Find(&Categories).RowsAffected
	if check == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Data Category",
		})
	}

	// update response struct
	var response []UpdateCategoryResponse
	for _, Category := range Categories {
		response = append(response, UpdateCategoryResponse{
			ID:        Category.ID,
			Name:      Category.Name,
			CreatedAt: Category.Created_at,
			UpdatedAt: Category.Updated_at,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    response,
	})
}

func GetCategory(c *fiber.Ctx) error {
	id := c.Query("id")
	var Category models.Category

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

	response := UpdateCategoryResponse{
		ID:        Category.ID,
		Name:      Category.Name,
		CreatedAt: Category.Created_at,
		UpdatedAt: Category.Updated_at,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    response,
	})

}

func UpdateCategory(c *fiber.Ctx) error {
	var cat models.Category
	if err := c.BodyParser(&cat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	id := c.Params("id")

	// Cek validasi data kategori
	if err := helper.ValidationStruct(c, cat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"error":   err,
		})
	}

	// Lakukan pembaruan data kategori ke database
	result := database.Database.Db.Model(&models.Category{}).Where("id = ?", id).Updates(&cat)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Ambil kategori yang berhasil diperbarui dari database
	var updatedCat models.Category
	if err := database.Database.Db.Where("id = ?", id).First(&updatedCat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Buat respons dengan data kategori yang berhasil diperbarui
	response := UpdateCategoryResponse{
		ID:        updatedCat.ID,
		Name:      updatedCat.Name,
		CreatedAt: updatedCat.Created_at,
		UpdatedAt: updatedCat.Updated_at,
	}

	// Kirim respons dengan data kategori yang berhasil diperbarui
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Update Category",
		"data":    response,
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
