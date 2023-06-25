package handler

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Flight struct {
	DB *gorm.DB
}

func (f *Flight) Get(c echo.Context) error {
	return nil
}
