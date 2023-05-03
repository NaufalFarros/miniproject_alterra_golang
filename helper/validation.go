package helper

import (
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
	// cek value yang dikrim dari client
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
