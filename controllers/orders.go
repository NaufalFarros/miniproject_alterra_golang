package controllers

import (
	"fmt"
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/gofiber/fiber/v2"
)

type Orders struct {
	ID            uint      `json:"id"`
	Name_customer string    `json:"name_customer" validate:"required"`
	Phone         string    `json:"phone" validate:"required"`
	Table_number  string    `json:"table_number" validate:"required"`
	UserID        int       `json:"user_id" validate:"required"`
	Created_at    time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItems struct {
	ID             uint      `json:"id"`
	ItemID         int       `json:"item_id" validate:"required"`
	UserID         int       `json:"user_id" validate:"required"`
	Quantity       int       `json:"quantity" validate:"required"`
	SubTotal       int       `json:"sub_total" validate:"required"`
	Quantity_total int       `json:"quantity_total" validate:"required"`
	Total_price    int       `json:"total_price" validate:"required"`
	Status_order   string    `json:"status_order" default:"pending" validate:"required"`
	OrdersID       int       `json:"orders_id" validate:"required"`
	Created_at     time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Table struct {
	Table_number string    `json:"table_number" validate:"required"`
	Created_at   time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func CreateBookings(c *fiber.Ctx) error {
	var booking Orders
	// var dbUsers = Users{}
	// Get user ID from JWT token
	userID := c.Locals("userID").(int)

	if err := c.BodyParser(&booking); err != nil {
		return err
	}
	var dataTable = Table{}
	if getID := database.Database.Db.Select("tables.table_number").Joins("JOIN users ON users.table_id = tables.id").Where("users.id = ?", userID).First(&dataTable); getID.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Table Not Found",
		})
	}

	booking.UserID = userID
	booking.Table_number = dataTable.Table_number

	errors := helper.ValidationStruct(c, booking)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	err := database.Database.Db.Create(&booking)

	if err.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Create Booking",
		"data":    booking,
	})

}

func CreateBookingsItems(c *fiber.Ctx) error {

	// Get user ID from JWT token
	userID := c.Locals("userID").(int)
	fmt.Println(userID)

	var orders OrderItems

	if err := c.BodyParser(&orders); err != nil {
		return err
	}

	var Items Items

	if err := database.Database.Db.Where("id = ?", orders.ItemID).First(&Items).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item Not Found",
		})
	}
	orders.UserID = int(userID)
	fmt.Println("ini user ID", orders.UserID)
	orders.Quantity = orders.Quantity + 1
	orders.SubTotal = Items.Price * orders.Quantity
	orders.Quantity_total = orders.Quantity
	orders.Total_price = orders.SubTotal

	var ordersID Orders

	if err := database.Database.Db.Where("user_id = ?", userID).First(&ordersID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Orders Not Found",
		})
	}

	orders.OrdersID = int(ordersID.ID)
	fmt.Println("ini OrDErs ID", orders.OrdersID)
	orders.Created_at = time.Now()
	orders.Updated_at = time.Now()

	var ordersItems OrderItems
	err := database.Database.Db.Where("orders_id = ? AND status_order = ? AND item_id = ?", orders.OrdersID, "pending", orders.ItemID).First(&ordersItems).Error

	fmt.Println("ini ordersItems", ordersItems)

	if err == nil {
		// Jika item sudah ada, tambahkan kuantitas dan subtotal baru
		ordersItems.Quantity += orders.Quantity
		ordersItems.SubTotal = Items.Price * ordersItems.Quantity
		ordersItems.Quantity_total += orders.Quantity
		ordersItems.Total_price = ordersItems.SubTotal
		database.Database.Db.Save(&ordersItems).Where("id = ?", ordersItems.ID)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Success Add Orders Items",
			"data":    ordersItems,
		})
	} else {
		// Jika item belum ada, buat baris baru di tabel orders_items
		newOrdersItem := OrderItems{
			UserID:         orders.UserID,
			ItemID:         orders.ItemID,
			OrdersID:       orders.OrdersID,
			Quantity:       orders.Quantity,
			SubTotal:       Items.Price * orders.Quantity,
			Quantity_total: orders.Quantity,
			Total_price:    Items.Price * orders.Quantity,
			Status_order:   "pending",
			Created_at:     time.Now(),
			Updated_at:     time.Now(),
		}
		database.Database.Db.Create(&newOrdersItem)
		// return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		// 	"message": "Success Add Orders Items",
		// 	"data":    orders,
		// })
	}
	errors := helper.ValidationStruct(c, orders)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Add Orders Items",
		"data":    orders,
	})
}

func CreateBookingsItemsMin(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	// fmt.Println(userID)

	var orders OrderItems

	if err := c.BodyParser(&orders); err != nil {
		return err
	}

	if err := database.Database.Db.Where("item_id = ? AND user_id = ? AND status_order =?", orders.ItemID, userID, "pending").First(&orders).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item Not Found or Status Order Not Pending",
		})
	}
	var Items Items
	if err := database.Database.Db.Where("id = ?", orders.ItemID).First(&Items).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item Not Found",
		})
	}
	orders.UserID = int(userID)
	orders.Quantity = orders.Quantity - 1
	orders.SubTotal = Items.Price * orders.Quantity
	orders.Quantity_total = orders.Quantity
	orders.Total_price = orders.SubTotal

	database.Database.Db.Save(&orders).Where("id = ?", orders.ID)

	if orders.Quantity == 0 {
		database.Database.Db.Delete(&orders).Where("id = ?", orders.ID)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success Delete Orders Items",
			// "data":    orders,
		})
	}

	errors := helper.ValidationStruct(c, orders)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Minus Orders Items",
		"data":    orders,
	})

}

func SubmitOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	var ordersItems OrderItems

	if err := c.BodyParser(&ordersItems); err != nil {
		return err
	}
	var orders Orders

	if err1 := database.Database.Db.Where("user_id = ? AND id =?", userID, ordersItems.OrdersID).First(&orders).Error; err1 != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Orders Not Found",
		})
	}

	ordersItems.UserID = int(userID)
	ordersItems.OrdersID = int(orders.ID)
	if err := database.Database.Db.Where("user_id = ? AND status_order =? AND order_id =?", userID, "pending").First(&ordersItems).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Orders Not Found or Status Order Not Pending",
		})
	}
	// errors := helper.ValidationStruct(c, ordersItems)

	// if errors != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(errors)
	// }

	database.Database.Db.Model(&OrderItems{}).Where("orders_id = ? AND status_order = ?", orders.ID, "pending").Updates(map[string]interface{}{
		"status_order": "paying",
		"updated_at":   time.Now(),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Submit Orders Items",
		"data":    ordersItems,
	})
}
