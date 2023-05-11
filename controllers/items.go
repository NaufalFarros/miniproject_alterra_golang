package controllers

import (
	"os"
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

type ItemResponse struct {
	ID                     uint                   `json:"id"`
	Name                   string                 `json:"name" validate:"required"`
	Image                  string                 `json:"image" validate:"required"`
	Description            string                 `json:"description" validate:"required"`
	Price                  int                    `json:"price" validate:"required"`
	Stock                  int                    `json:"stock" validate:"required"`
	UpdateCategoryResponse UpdateCategoryResponse `json:"category"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type ItemsResponse struct {
	Data    []ItemResponse `json:"data"`
	Message string         `json:"message"`
}

func itemToResponse(item models.Items) ItemResponse {
	itemResponse := ItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Image:       item.Image,
		Description: item.Description,
		Price:       item.Price,
		Stock:       item.Stock,
		UpdateCategoryResponse: UpdateCategoryResponse{
			ID:        item.Category.ID,
			Name:      item.Category.Name,
			CreatedAt: item.Category.Created_at,
			UpdatedAt: item.Category.Updated_at,
		},
		CreatedAt: item.Created_at,
		UpdatedAt: item.Updated_at,
	}

	return itemResponse
}

func CreateItem(c *fiber.Ctx) error {
	var item models.Items

	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	file, err := c.FormFile("Image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request image",
		})
	}

	var category models.Category
	if err := database.Database.Db.Where("id = ?", item.CategoryID).First(&category).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Category ID not found",
		})
	}

	file.Filename = helper.GenerateFileName(file.Filename)

	if err := c.SaveFile(file, "./image/"+file.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error at SaveFile",
		})
	}

	item.Name = c.FormValue("Name")
	item.Image = file.Filename
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	errors := helper.ValidationStruct(c, item)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	result := database.Database.Db.Create(&item)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Ambil kategori terkait dari database
	if err := database.Database.Db.Model(&item).Association("Category").Find(&item.Category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	response := itemToResponse(item)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Create items",
		"data":    response,
	})
}

func GetItems(c *fiber.Ctx) error {
	var items []models.Items
	if err := database.Database.Db.Preload("Category").Find(&items).Error; err != nil {
		// Tangani kesalahan saat mengambil data item
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get items",
			"error":   err.Error(),
		})
	}
	// cek jika data masih tidak ada
	check := database.Database.Db.Preload("Category").Find(&items).RowsAffected
	if check == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	for i := range items {
		items[i].Image = c.BaseURL() + "/images/" + items[i].Image
	}

	var itemsResponse []ItemResponse
	for _, item := range items {
		itemResponse := itemToResponse(item)
		itemsResponse = append(itemsResponse, itemResponse)
	}

	// Kirim respons dengan data item yang berhasil diambil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Get Items",
		"data":    itemsResponse,
	})
}

func GetItem(c *fiber.Ctx) error {
	var Item = []models.Items{}
	id := c.Query("id")
	// fmt.Println(id)

	errors := database.Database.Db.Preload("Category").Where("id = ?", id).First(&Item)

	if errors.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "No Data Found",
		})
	}

	// Tambahkan URL gambar ke objek JSON yang dikirim sebagai respons
	for i := range Item {
		Item[i].Image = c.BaseURL() + "/images/" + Item[i].Image
	}

	var itemsResponse []ItemResponse
	for _, items := range Item {
		itemResponse := itemToResponse(items)
		itemsResponse = append(itemsResponse, itemResponse)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Get Items",
		"data":    itemsResponse,
	})
}

func UpdateItem(c *fiber.Ctx) error {
	var item = models.Items{}

	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	id := c.Params("id")

	updates := make(map[string]interface{})

	if item.Name != "" {
		updates["name"] = item.Name
	}

	// Pengecekan apakah ada file gambar baru
	file, err := c.FormFile("Image")
	if err == nil {
		file.Filename = helper.GenerateFileName(file.Filename)
		checkimage := database.Database.Db.Where("id = ?", id).First(&item)
		if checkimage.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "iamge not found",
			})
		}
		if err := os.Remove("./image/" + item.Image); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error at Remove old Image",
			})
		}

		if err := c.SaveFile(file, "./image/"+file.Filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error at SaveFile",
			})
		}
		updates["image"] = file.Filename
	}

	if item.Description != "" {
		updates["description"] = item.Description
	}

	if item.Price != 0 {
		updates["price"] = item.Price
	}

	if item.Stock != 0 {
		updates["stock"] = item.Stock
	}

	if item.CategoryID != 0 {
		updates["category_id"] = item.CategoryID
	}

	updates["updated_at"] = time.Now()

	// validasi

	checkID := database.Database.Db.Where("id = ?", id).First(&item)
	if checkID.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID not found",
		})
	}

	result := database.Database.Db.Model(&models.Items{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	database.Database.Db.Preload("Category").Where("id = ?", id).Find(&item)

	item.Image = c.BaseURL() + "/images/" + item.Image

	// update response
	response := itemToResponse(item)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Update Item",
		"data":    response,
	})
}

func DeleteItem(c *fiber.Ctx) error {
	var Item = models.Items{}

	checkID := database.Database.Db.Where("id = ?", c.Params("id")).First(&Item)
	if checkID.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID not found",
		})
	}

	if err := os.Remove("./image/" + Item.Image); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error at Remove old Image",
		})
	}

	errors := database.Database.Db.Where("id = ?", c.Params("id")).Delete(&Item)

	if errors.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Delete Item",
	})
}
