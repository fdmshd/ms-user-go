package utils

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func (cv *CustomValidator) ValidateExcept(i interface{}, fields ...string) error {
	if err := cv.validator.StructExcept(i, fields...); err != nil {
		return err
	}
	return nil
}

func NewValidator() *CustomValidator {
	v := validator.New()
	return &CustomValidator{validator: v}
}
