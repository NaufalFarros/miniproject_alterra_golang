package helper

import (
	"math/rand"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidationStruct(c *fiber.Ctx, model interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ErrorResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			})
		}
		return errors
	}
	return nil
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateFileName(fileName string) string {
	return "item-" + RandomString(6) + filepath.Ext(fileName)
}
