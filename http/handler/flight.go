package handler

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

type GetFlightsRequest struct {
	DepCity string     `query:"departure_city" validate:"required"`
	ArrCity string     `query:"arrival_city" validate:"required"`
	DepTime *time.Time `query:"date" validate:"required"`
}

type City struct {
	ID   int32
	Name string
}

type Airplane struct {
	ID   int32
	Name string
}

type GetFlightsResponse struct {
	ID             int32     `json:"id"`
	DepCity        City      `json:"dep_city"`
	ArrCity        City      `json:"arr_city"`
	DepTime        time.Time `json:"dep_time"`
	ArrTime        time.Time `json:"arr_time"`
	Airplane       Airplane  `json:"airplane"`
	Airline        string    `json:"airline"`
	Price          int32     `json:"price"`
	CxlSitID       int32     `json:"cxl_sit_id"`
	RemainingSeats int32     `json:"remaining_seats"`
}

func (f *Flight) Get(c echo.Context) error {
	var req GetFlightsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	if err := f.Validator.Struct(&req); err != nil {
		return c.JSONPretty(http.StatusBadRequest, "Bad Request", " ")
	}

	var flights []models.Flight
	err := f.DB.Debug().
		Joins("Airplane").
		Joins("DepCity").Where("DepCity.name = ?", req.DepCity).
		Joins("ArrCity").Where("ArrCity.name = ?", req.ArrCity).
		Where("year(dep_time) = ?", req.DepTime.Year()).
		Where("month(dep_time) = ?", req.DepTime.Month()).
		Where("day(dep_time) = ?", req.DepTime.Day()).
		Find(&flights).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	response := make([]GetFlightsResponse, 0, len(flights))
	for _, val := range flights {
		response = append(response, GetFlightsResponse{
			ID:             val.ID,
			DepCity:        City{ID: val.DepCity.ID, Name: val.DepCity.Name},
			ArrCity:        City{ID: val.ArrCity.ID, Name: val.ArrCity.Name},
			DepTime:        val.DepTime,
			ArrTime:        val.ArrTime,
			Airplane:       Airplane{ID: val.Airplane.ID, Name: val.Airplane.Name},
			Airline:        val.Airline,
			Price:          val.Price,
			CxlSitID:       val.CxlSitID,
			RemainingSeats: val.RemainingSeats,
		})
	}

	return c.JSONPretty(http.StatusOK, response, " ")
}
