package handler

import (
	"flight-data-api/models"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FlightsDates struct {
	DB *gorm.DB
}

type GetFlightsDatesResponse struct {
	Dates []string
}

func (f *FlightsDates) Get(ctx echo.Context) error {
	var dates []time.Time
	err := f.DB.Debug().
		Model(&models.Flight{}).
		Distinct("DATE(dep_time)").
		Scan(&dates).
		Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	var response GetFlightsDatesResponse
	response.Dates = make([]string, 0, len(dates))
	for _, date := range dates {
		response.Dates = append(response.Dates, date.Format("2006-01-02"))
	}

	sort.Slice(response.Dates, func(i, j int) bool {
		return response.Dates[i] < response.Dates[j]
	})
	return ctx.JSONPretty(http.StatusOK, response, " ")
}
