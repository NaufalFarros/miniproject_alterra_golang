package controllers

import (
	"github.com/NaufalFarros/miniproject_alterra_golang/config"
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID                   uint                 `json:"id"`
	Email                string               `json:"email"`
	Name                 string               `json:"name"`
	CustomeResponseTable CustomeResponseTable `json:"table"`
	CustomeRoles         CustomeRoles         `json:"role"`
}

type UsersResponse struct {
	Data    []UserResponse `json:"data"`
	Message string         `json:"message"`
}

type CustomAllOrdersResponse struct {
	ID           uint                 `json:"id"`
	NameCustomer string               `json:"name_customer"`
	Phone        string               `json:"phone"`
	Status_Order string               `json:"status_order"`
	Users        CustomeResponseUsers `json:"users"`
	Order        []CustomeOrdersItems `json:"orders_items"`
}

type CustomeResponseUsers struct {
	ID                   uint                 `json:"id"`
	Name                 string               `json:"name"`
	Email                string               `json:"email"`
	CustomeResponseTable CustomeResponseTable `json:"table"`
	CustomeRoles         CustomeRoles         `json:"roles"`
}

type CustomeResponseTable struct {
	ID          uint   `json:"id"`
	TableNumber string `json:"table_number"`
}

type CustomeOrdersItems struct {
	ID            uint `json:"id"`
	Quantity      int  `json:"quantity"`
	SubTotal      int  `json:"sub_total"`
	QuantityTotal int  `json:"quantity_total"`
	TotalPrice    int  `json:"total_price"`
	OrdersID      uint `json:"orders_id"`
	CustomeItems  `json:"items"`
}

type CustomeItems struct {
	ID              uint            `json:"id"`
	Name            string          `json:"name"`
	Price           int             `json:"price"`
	Stock           int             `json:"stock"`
	Image           string          `json:"image"`
	CustomeCategory CustomeCategory `json:"category"`
}

type CustomeCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CustomeRoles struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	TableID  int    `json:"table_id" validate:"required"`
	RoleID   int    `json:"role_id" validate:"required"`
}

func userToResponse(user models.User) UserResponse {
	userResponse := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		CustomeResponseTable: CustomeResponseTable{
			ID:          user.Table.ID,
			TableNumber: user.Table.Table_number,
		},
		CustomeRoles: CustomeRoles{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		},
	}
	return userResponse
}

func HashingPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	var users = models.User{}
	if err := c.BodyParser(&users); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	result := database.Database.Db.Preload("Role").Preload("Table").Where("email = ?", users.Email).First(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
		})
	}

	var role string
	database.Database.Db.Table("roles").Select("roles.name").Joins("JOIN users ON users.role_id = roles.id").Where("users.id = ?", users.ID).Scan(&role)

	if !CheckPasswordHash(c.FormValue("password"), users.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Password or Email",
		})
	}

	token, err := config.CreateToken(users.Email, role, int(users.ID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error Check Token",
		})
	}
	userResponse := []UserResponse{}
	userResponse = append(userResponse, userToResponse(users))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Login",
		"token":   token,
		"data":    userResponse,
	})
}

func Register(c *fiber.Ctx) error {
	var user = models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	checkEmail := database.Database.Db.Where("email = ?", user.Email).First(&user)

	if checkEmail.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email Already Registered",
		})
	}
	errors := helper.ValidationStruct(c, user)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}
	hash, err := HashingPassword(user.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	user.Password = hash
	user.RoleID = 2

	result := database.Database.Db.Create(&user)
	database.Database.Db.Preload("Role").Preload("Table").Where("id = ?", user.ID).First(&user)
	userResponse := []UserResponse{}
	userResponse = append(userResponse, userToResponse(user))

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error DB",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success",
		"data":    userResponse,
	})
}

func Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Logout",
	})
}

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	result := database.Database.Db.Preload("Role").Preload("Table").Find(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	userResponse := []UserResponse{}
	for _, user := range users {
		usersResponse := userToResponse(user)
		userResponse = append(userResponse, usersResponse)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    userResponse,
	})
}

func GetAllOrdersUsers(c *fiber.Ctx) error {
	var orders []models.Orders
	var customAllOrdersResponse []CustomAllOrdersResponse

	err := database.Database.Db.
		Preload("User").
		Preload("User.Table").
		Preload("User.Role").
		Preload("OrderItems").
		Preload("OrderItems.Item").
		Preload("OrderItems.Item.Category").
		Find(&orders).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error",
			"error":   err.Error(),
		})
	}

	for _, order := range orders {
		var orderItems []models.OrderItems
		err := database.Database.Db.
			Where("orders_id = ?", order.ID).
			Preload("Item").
			Preload("Item.Category").
			Find(&orderItems).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error",
				"error":   err.Error(),
			})
		}

		var customOrderItems []CustomeOrdersItems
		for _, orderItem := range orderItems {
			customOrderItem := CustomeOrdersItems{
				ID:            orderItem.ID,
				Quantity:      orderItem.Quantity,
				SubTotal:      orderItem.SubTotal,
				QuantityTotal: orderItem.Quantity_total,
				TotalPrice:    orderItem.Total_price,
				OrdersID:      orderItem.OrdersID,
				CustomeItems: CustomeItems{
					ID:    orderItem.Item.ID,
					Name:  orderItem.Item.Name,
					Price: orderItem.Item.Price,
					Stock: orderItem.Item.Stock,
					Image: orderItem.Item.Image,
					CustomeCategory: CustomeCategory{
						ID:   orderItem.Item.Category.ID,
						Name: orderItem.Item.Category.Name,
					},
				},
			}
			customOrderItems = append(customOrderItems, customOrderItem)
		}

		customOrder := CustomAllOrdersResponse{
			ID:           order.ID,
			NameCustomer: order.Name_customer,
			Phone:        order.Phone,
			Status_Order: order.Status_order,
			Users: CustomeResponseUsers{
				ID:    order.User.ID,
				Name:  order.User.Name,
				Email: order.User.Email,
				CustomeResponseTable: CustomeResponseTable{
					ID:          order.User.Table.ID,
					TableNumber: order.User.Table.Table_number,
				},
				CustomeRoles: CustomeRoles{
					ID:   order.User.Role.ID,
					Name: order.User.Role.Name,
				},
			},
			Order: customOrderItems,
		}

		customAllOrdersResponse = append(customAllOrdersResponse, customOrder)
	}

	// cek apakah ada data
	if len(customAllOrdersResponse) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    customAllOrdersResponse,
	})
}

func UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var order Orders
	if err := database.Database.Db.First(&order, orderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	order.Status_order = "payment_success"
	if err := database.Database.Db.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update order status",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order status updated successfully",
		"data":    order,
	})
}
