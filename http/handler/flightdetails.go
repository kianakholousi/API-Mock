package handler

import (
	"API-Mock/models"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type FlightDetail struct {
	DB *gorm.DB
}

func (f *FlightDetail) Get(c echo.Context) error {
	// Convert param (string) to int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	// Initialize struct
	var flight models.Flight

	// Find flight by ID, select only FlightClass field
	result := f.DB.Select("flight_class").First(&flight, id)

	// If record not found
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Flight not found")
	} else if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching flight class")
	}

	// If all went well, return the FlightClass
	return c.JSON(http.StatusOK, flight.FlightClass)
}

