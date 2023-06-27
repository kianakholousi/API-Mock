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

type FlightGetRequest struct {
	DepCity    string `query:"departure_city"`
	ArrCity    string `query:"arrival_city"`
	DepTimeStr string `query:"departure_time"`
	//DepTime    time.Time `query:"-"`
}

func (f *Flight) Get(c echo.Context) error {
	var FGR FlightGetRequest
	if err := c.Bind(&FGR); err != nil {
		return err
	}

	DepTime, err := time.Parse("2006-01-02", FGR.DepTimeStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	//if arrCity == "" || depCity == "" {
	//	return c.JSON(http.StatusBadRequest, "Bad Request")
	//}

	var flights []models.Flight
	err = f.DB.Debug().Joins("DepCity").Where("DepCity.name = ?", FGR.DepCity).
		Joins("ArrCity").Where("ArrCity.name = ?", FGR.ArrCity).
		Where("year(dep_time) = ?", DepTime.Year()).
		Where("month(dep_time) = ?", DepTime.Month()).
		Where("day(dep_time) = ?", DepTime.Day()).
		Find(&flights).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSONPretty(http.StatusOK, flights, " ")
}
