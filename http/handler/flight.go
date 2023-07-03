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

type GetFlightsResponse struct {
	ID             int32                `json:"id"`
	DepCity        GetCitiesResponse    `json:"dep_city"`
	ArrCity        GetCitiesResponse    `json:"arr_city"`
	DepTime        time.Time            `json:"dep_time"`
	ArrTime        time.Time            `json:"arr_time"`
	Airplane       GetAirplanesResponse `json:"airplane"`
	Airline        string               `json:"airline"`
	Price          int32                `json:"price"`
	CxlSitID       int32                `json:"cxl_sit_id"`
	RemainingSeats int32                `json:"remaining_seats"`
}

func (f *Flight) Get(ctx echo.Context) error {
	var req GetFlightsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	}

	if err := f.Validator.Struct(&req); err != nil {
		return ctx.JSONPretty(http.StatusBadRequest, "Bad Request", " ")
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
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	response := make([]GetFlightsResponse, 0, len(flights))
	for _, val := range flights {
		response = append(response, GetFlightsResponse{
			ID:             val.ID,
			DepCity:        GetCitiesResponse{ID: val.DepCity.ID, Name: val.DepCity.Name},
			ArrCity:        GetCitiesResponse{ID: val.ArrCity.ID, Name: val.ArrCity.Name},
			DepTime:        val.DepTime,
			ArrTime:        val.ArrTime,
			Airplane:       GetAirplanesResponse{ID: val.Airplane.ID, Name: val.Airplane.Name},
			Airline:        val.Airline,
			Price:          val.Price,
			CxlSitID:       val.CxlSitID,
			RemainingSeats: val.RemainingSeats,
		})
	}

	return ctx.JSONPretty(http.StatusOK, response, " ")
}
