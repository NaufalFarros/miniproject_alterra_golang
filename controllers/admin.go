package controllers

import (
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/gofiber/fiber/v2"
)

type CustomAllOrdersResponse struct {
	ID                   uint                 `json:"id"`
	NameCustomer         string               `json:"name_customer"`
	Phone                string               `json:"phone"`
	TableNumber          int                  `json:"table_number"`
	Status_Order         string               `json:"status_order"`
	UsersId              uint                 `json:"users_id"`
	CustomeResponseUsers CustomeResponseUsers `json:"custome_response_users"`
	CustomeResponseTable CustomeResponseTable `json:"custome_response_table"`
	CustomeRoles         CustomeRoles         `json:"custome_roles"`
	CustomeUsersOrders   CustomeUsersOrders   `json:"custome_users_orders"`
	CustomeOrdersItems   CustomeOrdersItems   `json:"custome_orders_items"`
	CustomeCategory      CustomeCategory      `json:"custome_category"`
}

type CustomeResponseUsers struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	TableID int    `json:"table_id"`
}

type CustomeResponseTable struct {
	ID          uint `json:"id"`
	TableNumber int  `json:"table_number"`
}

type CustomeUsersOrders struct {
	ID            uint `json:"id"`
	ItemID        int  `json:"item_id"`
	UserID        int  `json:"user_id"`
	Quantity      int  `json:"quantity"`
	SubTotal      int  `json:"sub_total"`
	QuantityTotal int  `json:"quantity_total"`
	TotalPrice    int  `json:"total_price"`
	OrdersID      int  `json:"orders_id"`
}

type CustomeOrdersItems struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	Image      string `json:"image"`
	CategoryID int    `json:"category_id"`
}

type CustomeCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CustomeRoles struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetAllOrdersUsers(c *fiber.Ctx) error {
	var customAllOrdersResponse []CustomAllOrdersResponse

	// SQL query
	err := database.Database.Db.Table("order_items").
		Select(`
			order_items.id,
			users.name AS name_customer,
			users.phone,
			tables.number AS table_number,
			orders.status AS status_order,
			users.id AS users_id,
			orders.id AS orders_id,
			orders.id AS orders_id,
			orders.quantity AS quantity,
			orders.sub_total AS sub_total,
			orders.quantity_total AS quantity_total,
			orders.total_price AS total_price,
			items.id AS item_id,
			items.name AS item_name,
			items.price AS item_price,
			items.stock AS item_stock,
			items.image AS item_image,
			categories.id AS category_id,
			categories.name AS category_name,
			user_roles.id AS role_id,
			user_roles.name AS role_name
		`).
		Joins(`
			JOIN orders ON orders.id = order_items.orders_id
			JOIN users ON users.id = orders.user_id
			JOIN tables ON tables.table_number = orders.table_number
			JOIN roles AS user_roles ON roles.id = users.role_id
			JOIN order_items ON order_items.order_id = orders.id
			JOIN items ON items.id = order_items.item_id
			JOIN categories ON categories.id = items.category_id
		`).
		Scan(&customAllOrdersResponse).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    customAllOrdersResponse,
	})
}
