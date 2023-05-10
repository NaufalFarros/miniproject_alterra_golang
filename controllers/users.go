package controllers

import (
	"github.com/NaufalFarros/miniproject_alterra_golang/config"
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)


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

