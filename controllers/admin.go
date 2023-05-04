package controllers

import (
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/config"
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/helper"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// hash password
func HashingPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// compare password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	var Users = models.Admin{}
	Users.Email = c.FormValue("email")
	Users.Password = c.FormValue("password")

	// find user by username
	result := database.Database.Db.Where("email = ?", Users.Email).First(&Users)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
		})
	}

	// compare password
	if !CheckPasswordHash(c.FormValue("password"), Users.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Password or Email",
		})
	}

	// generate token
	token, err := config.CreateToken(c.FormValue("email"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// enctkn, err := config.EncryptToken(token)

	// if err != nil {
	// 	log.Println(err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"message": "Internal Server Error ecttkn",
	// 	})
	// }

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Login",
		"token":   token,
	})

}

func Register(c *fiber.Ctx) error {
	var Users = models.Admin{}

	if err := c.BodyParser(&Users); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}
	// cek validasi
	errors := helper.ValidationStruct(c, Users)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// hash password
	hash, err := HashingPassword(Users.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	Users.Password = hash
	Users.CreatedAt = time.Now()
	Users.UpdatedAt = time.Now()

	// save to database

	result := database.Database.Db.Create(&Users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success",
	})

}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    "Authorization",
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour * 24),
	})

	// remove token jwt
	// token := c.Locals("user").(*jwt.Token)
	// token.Valid = false

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Logout",
	})
}

func GetUsers(c *fiber.Ctx) error {
	var Users []models.Admin

	result := database.Database.Db.Find(&Users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    Users,
	})
}

func GetUser(c *fiber.Ctx) error {
	var Users []models.Admin

	result := database.Database.Db.Where("id = ?", c.Params("id")).First(&Users)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    Users,
	})

}




