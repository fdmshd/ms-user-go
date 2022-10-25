package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= 6
	})
	v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= 3
	})
	return &CustomValidator{validator: v}
}
