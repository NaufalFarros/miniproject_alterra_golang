package controllers

import (
	"time"

	"github.com/NaufalFarros/miniproject_alterra_golang/config"
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID        uint         `json:"id"`
	Email     string       `json:"email"`
	Name      string       `json:"name"`
	TableID   int          `json:"table_id"`
	Table     models.Table `json:"table"`
	RoleID    int          `json:"role_id"`
	Role      models.Roles `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type UsersResponse struct {
	Data    []UserResponse `json:"data"`
	Message string         `json:"message"`
}

func userToResponse(user models.User) UserResponse {
	userResponse := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		TableID:   user.TableID,
		Table:     user.Table,
		RoleID:    user.RoleID,
		Role:      user.Role,
		CreatedAt: user.Created_at,
		UpdatedAt: user.Updated_at,
	}
	return userResponse
}

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
	var users = models.User{}
	if err := c.BodyParser(&users); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	// find user by Email
	result := database.Database.Db.Preload("Role").Preload("Table").Where("email = ?", users.Email).First(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User Not Found",
		})
	}

	var role string
	database.Database.Db.Table("roles").Select("roles.name").Joins("JOIN users ON users.role_id = roles.id").Where("users.id = ?", users.ID).Scan(&role)

	// compare password
	if !CheckPasswordHash(c.FormValue("password"), users.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Password or Email",
		})
	}

	// generate token
	token, err := config.CreateToken(users.Email, role, int(users.ID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error Check Token",
		})
	}
	userResponse := []UserResponse{}
	userResponse = append(userResponse, userToResponse(users))

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

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

	//validasi jika email sudah terdaftar di database
	checkEmail := database.Database.Db.Where("email = ?", user.Email).First(&user)

	if checkEmail.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email Already Registered",
		})
	}

	// hash password
	hash, err := HashingPassword(user.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	user.Password = hash
	user.RoleID = 2
	user.Created_at = time.Now()
	user.Updated_at = time.Now()

	// cek validasi
	// errors := helper.ValidationStruct(c, user)
	// if errors != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(errors)
	// }

	// save to database
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

func GetUser(c *fiber.Ctx) error {
	var Users []models.User

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
