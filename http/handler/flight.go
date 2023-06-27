package handler

import "C"
import (
	"flight-data-api/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
	"net/http"
	"time"
)

type Flight struct {
	DB        *gorm.DB
	Validator *validator.Validate
}

type FlightsGetRequest struct {
	DepCity string     `query:"departure_city" validate:"required"`
	ArrCity string     `query:"arrival_city" validate:"required"`
	DepTime *time.Time `query:"departure_time" validate:"required"`
}

func (f *Flight) Get(c echo.Context) error {
	var req FlightsGetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSONPretty(http.StatusBadRequest, "Bad Request", " ")
	}

	if err := f.Validator.Struct(&req); err != nil {
		return c.JSONPretty(http.StatusBadRequest, "Bad Request", " ")
	}

	var flights []models.Flight
	err := f.DB.Debug().Joins("DepCity").
		Where("DepCity.name = ?", req.DepCity).
		Joins("ArrCity").Where("ArrCity.name = ?", req.ArrCity).
		Where("year(dep_time) = ?", req.DepTime.Year()).
		Where("month(dep_time) = ?", req.DepTime.Month()).
		Where("day(dep_time) = ?", req.DepTime.Day()).
		Find(&flights).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSONPretty(http.StatusOK, flights, " ")
}
