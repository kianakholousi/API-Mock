package handler

import (
	"flight-data-api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Flight struct {
	DB *gorm.DB
}

func (f *Flight) FlightGet(c echo.Context) error {
	var ff []models.Flight
	//err := f.DB.Preload(clause.Associations).Where("dep_city_id = 2").First(&ff).Error
	err := f.DB.Model(&models.Flight{}).Joins("ArrCity", models.City{Name: "Mashhad"}).Find(&ff).Error
	if err != nil {
		panic(err)
	}
	if err := c.JSONPretty(http.StatusOK, ff, "	"); err != nil {
		panic(err)
	}
	return nil
}
