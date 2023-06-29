package handler

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Airplanes struct {
	DB *gorm.DB
}

func (a *Airplanes) Get(ctx echo.Context) error {
	return nil
}
