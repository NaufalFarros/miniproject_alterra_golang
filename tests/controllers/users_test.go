package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NaufalFarros/miniproject_alterra_golang/controllers"
	"github.com/NaufalFarros/miniproject_alterra_golang/database"
	"github.com/NaufalFarros/miniproject_alterra_golang/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func TestRegister(t *testing.T) {
	database.DBConnect()
	app := fiber.New()

	app.Post("/register", controllers.Register)

	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		TableID:  1,
		RoleID:   1,
	}

	payload, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(body))
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %v but got %v", http.StatusCreated, resp.StatusCode)
	}

	database.DisconnectDB()
}

func TestLogin(t *testing.T) {
	err := godotenv.Load("./.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	database.DBConnect()
	app := fiber.New()

	app.Post("/login", controllers.Login)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("email", "test@example.com")
	_ = writer.WriteField("password", "password123")
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body = new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(body.String())

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v but got %v", http.StatusOK, resp.StatusCode)
	}

	database.DisconnectDB()
}

func TestLogout(t *testing.T) {
	app := fiber.New()
	app.Post("/logout", controllers.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v but got %v", http.StatusOK, resp.StatusCode)
	}
	database.DisconnectDB()
}
