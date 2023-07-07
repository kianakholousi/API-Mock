package handler

import (
	"errors"
	"flight-data-api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type FlightDetail struct {
	DB *gorm.DB
}

func (f *FlightDetail) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	var flight models.Flight
	result := f.DB.Model(&models.Flight{}).Select("id = ?", id).Find(&flight)

	// If record not found
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Flight not found")
	} else if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching flight class")
	}

	// If all went well, return the FlightClass
	return c.JSON(http.StatusOK, flight.FlightClass)
}
