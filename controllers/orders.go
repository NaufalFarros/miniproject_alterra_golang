package controllers

import (
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
)

type Orders struct {
	ID            uint      `json:"id"`
	Name_customer string    `json:"name_customer" validate:"required"`
	Phone         string    `json:"phone" validate:"required"`
	Status_order  string    `json:"status_order" validate:"required"`
	Table_number  string    `json:"table_number" validate:"required"`
	UserID        int       `json:"user_id"`
	Created_at    time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItems struct {
	ID             uint      `json:"id"`
	ItemID         int       `json:"item_id" validate:"required"`
	Quantity       int       `json:"quantity" validate:"required"`
	SubTotal       int       `json:"sub_total" validate:"required"`
	Quantity_total int       `json:"quantity_total" validate:"required"`
	Total_price    int       `json:"total_price" validate:"required"`
	OrdersID       int       `json:"orders_id" validate:"required"`
	Created_at     time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Table struct {
	Table_number string `json:"table_number" validate:"required"`
}

func CreateBookings(c *fiber.Ctx) error {
	var booking Orders
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
	booking.Status_order = "pending"

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

	userID := c.Locals("userID").(int)

	var orders OrderItems

	if err := c.BodyParser(&orders); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	var Items models.Items

	if err := database.Database.Db.Where("id = ?", orders.ItemID).First(&Items).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item Not Found",
		})
	}

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
	if ordersID.Status_order != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Order Status is Not Pending",
		})
	}

	orders.OrdersID = int(ordersID.ID)

	var ordersItems OrderItems
	err := database.Database.Db.
		Joins("JOIN orders ON order_items.orders_id = orders.id").
		Where("order_items.orders_id = ? AND order_items.item_id = ? AND orders.user_id = ? AND orders.status_order = ?", orders.OrdersID, orders.ItemID, userID, "pending").
		First(&ordersItems).Error

	if err == nil {

		ordersItems.Quantity += orders.Quantity
		ordersItems.SubTotal = Items.Price * ordersItems.Quantity
		ordersItems.Quantity_total += orders.Quantity
		ordersItems.Total_price = ordersItems.SubTotal

		check := helper.ValidationStruct(c, orders)
		if check != nil {
			return c.Status(fiber.StatusBadRequest).JSON(check)
		}
		database.Database.Db.Model(&ordersItems).Updates(ordersItems)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Success Add Orders Items",
			"data":    ordersItems,
		})
	}

	database.Database.Db.Create(&orders)

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

	if err := database.Database.Db.
		Joins("JOIN orders ON order_items.orders_id = orders.id").
		Where(" order_items.item_id = ? AND orders.user_id =? AND orders.status_order = ? ", orders.ItemID, userID, "pending").First(&orders).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item Not Found or Status Order Not Pending",
		})
	}
	var Items models.Items
	if err := database.Database.Db.Where("id = ?", orders.ItemID).First(&Items).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Item Not Found",
		})
	}
	orders.Quantity = orders.Quantity - 1
	orders.SubTotal = Items.Price * orders.Quantity
	orders.Quantity_total = orders.Quantity
	orders.Total_price = orders.SubTotal

	database.Database.Db.Save(&orders).Where("id = ?", orders.ID)

	if orders.Quantity == 0 {
		database.Database.Db.Delete(&orders).Where("id = ?", orders.ID)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success Delete Orders Items",
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

	var order Orders

	if err := database.Database.Db.Where("user_id = ? AND status_order = ?", userID, "pending").First(&order).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Orders Not Found or Status Order Not Pending",
		})
	}

	order.Status_order = "payment_pending"

	errors := helper.ValidationStruct(c, order)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var orderItems []OrderItems
	if err := database.Database.Db.Where("orders_id = ?", order.ID).Find(&orderItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve order items",
		})
	}

	for _, orderItem := range orderItems {

		var item models.Items
		if err := database.Database.Db.Where("id = ?", orderItem.ItemID).First(&item).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve item",
			})
		}
		item.Stock = item.Stock - orderItem.Quantity

		database.Database.Db.Model(&item).Where("id = ?", item.ID).Update("stock", item.Stock)

	}

	database.Database.Db.Save(&order).Where("id = ?", order.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Submit Orders Items",
		"data":    order,
	})
}
