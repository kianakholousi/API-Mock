package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Cities struct {
	DB        *gorm.DB
	Validator *validator.Validate
}

func (c *Cities) Get(ctx echo.Context) error {
	return nil
}
