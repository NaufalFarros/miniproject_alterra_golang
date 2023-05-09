package controllers

import (
	"fmt"
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/gofiber/fiber/v2"
)

type Items struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	Stock       int       `json:"stock" validate:"required"`
	CategoryID  int       `json:"category_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func CreateItem(c *fiber.Ctx) error {
	var items Items

	if err := c.BodyParser(&items); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}
	file, err := c.FormFile("Image") // ambil file dari form-data
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	// merubah filename menjadi unik item + extensi
	file.Filename = helper.GenerateFileName(file.Filename)

	// simpan file ke folder uploads
	if err := c.SaveFile(file, "./image/"+file.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error at SaveFile",
		})
	}
	//DEBUG data yang diinput
	fmt.Println(items)
	items.Name = c.FormValue("Name")
	items.Image = file.Filename
	// fmt.Println(items.Image)
	// fmt.Println(items.Name)
	// fmt.Println(items.Price)
	// fmt.Println(items.Stock)
	// fmt.Println(items.Description)
	// fmt.Println(items.CategoryID)

	items.CreatedAt = time.Now()
	items.UpdatedAt = time.Now()
	// cek validasi
	errors := helper.ValidationStruct(c, items)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// save to database
	result := database.Database.Db.Create(&items)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Create items",
		"data":    items,
	})
}

func GetItems(c *fiber.Ctx) error {
	var Items []Items

	database.Database.Db.Find(&Items)

	// Tambahkan URL gambar ke objek JSON yang dikirim sebagai respons
	for i := range Items {
		Items[i].Image = c.BaseURL() + "/admin/images/" + Items[i].Image
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Get Items",
		"data":    Items,
	})
}

func GetItem(c *fiber.Ctx) error {
	var Item []Items
	id := c.Query("id")
	// fmt.Println(id)

	errors := database.Database.Db.Where("id = ?", id).Find(&Item)

	if errors.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Tambahkan URL gambar ke objek JSON yang dikirim sebagai respons
	for i := range Item {
		Item[i].Image = c.BaseURL() + "/admin/images/" + Item[i].Image
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Get Items",
		"data":    Item,
	})
}

func UpdateItem(c *fiber.Ctx) error {
	var Item Items

	if err := c.BodyParser(&Item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	id := c.Query("id")

	updates := make(map[string]interface{})

	if Item.Name != "" {
		updates["name"] = Item.Name
	}

	if Item.Image != "" {
		file, err := c.FormFile("Image") // ambil file dari form-data
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Bad Request",
			})
		}

		// merubah filename menjadi unik item + extensi
		file.Filename = helper.GenerateFileName(file.Filename)

		// simpan file ke folder uploads
		if err := c.SaveFile(file, "./image/"+file.Filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error at SaveFile",
			})
		}
		updates["image"] = file.Filename
	}

	if Item.Description != "" {
		updates["description"] = Item.Description
	}

	if Item.Price != 0 {
		updates["price"] = Item.Price
	}

	if Item.Stock != 0 {
		updates["stock"] = Item.Stock
	}

	if Item.CategoryID != 0 {
		updates["category_id"] = Item.CategoryID
	}

	updates["updated_at"] = time.Now()

	result := database.Database.Db.Model(&Item).Where("id = ?", id).Updates(Item)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// mengambil data yang sudah diupdate
	database.Database.Db.Where("id = ?", id).Find(&Item)

	// Tambahkan URL gambar ke objek JSON yang dikirim sebagai respons
	Item.Image = c.BaseURL() + "/admin/images/" + Item.Image

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Update Items",
		"data":    Item,
	})

}

func DeleteItem(c *fiber.Ctx) error {
	var Item Items

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
