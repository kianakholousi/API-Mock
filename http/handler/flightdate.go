package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FlightsDays struct {
	DB *gorm.DB
}

type GetFlightsDaysResponse struct {
	Dates   []time.Time `query:"dates"`
	DepTime time.Time   `json:"dep_time"`
	ArrTime time.Time   `json:"arr_time"`
}

func (f *FlightsDays) Get(c echo.Context) error {
	var dates []time.Time
	err := f.DB.Debug().
		Raw(`SELECT DISTINCT DATE(dep_time, arr_time) FROM FlightsDays `).
		Scan(&dates).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	response := GetFlightsDaysResponse{
		Dates: dates,
	}
	return c.JSONPretty(http.StatusOK, response, " ")
}
