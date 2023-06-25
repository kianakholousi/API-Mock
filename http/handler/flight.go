package handler

import "C"
import (
	"flight-data-api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Flight struct {
	DB *gorm.DB
}

func (f *Flight) Get(c echo.Context) error {
	arrCity := c.FormValue("arr_city")
	depCity := c.FormValue("dep_city")
	depTime, err := time.Parse("2006-01-02", c.FormValue("dep_time"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	if arrCity == "" || depCity == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	var flights []models.Flight
	err = f.DB.Debug().Joins("DepCity").Where("DepCity.name = ?", depCity).
		Joins("ArrCity").Where("ArrCity.name = ?", arrCity).
		Where("year(dep_time) = ?", depTime.Year()).
		Where("month(dep_time) = ?", depTime.Month()).
		Where("day(dep_time) = ?", depTime.Day()).
		Find(&flights).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSONPretty(http.StatusOK, flights, " ")
}
