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

type FlightsGetResponse struct {
	ID             int32                     `json:"id"`
	DepCity        models.City               `json:"dep_city"`
	ArrCity        models.City               `json:"arr_city"`
	DepTime        time.Time                 `json:"dep_time"`
	ArrTime        time.Time                 `json:"arr_time"`
	Airplane       models.Airplane           `json:"-"`
	Airline        string                    `json:"airline"`
	Price          int32                     `json:"price"`
	CxlSit         models.CancelingSituation `json:"cxl_sit"`
	RemainingSeats int32                     `json:"remaining_seats"`
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

	response := make([]FlightsGetResponse, 0)
	for _, val := range flights {
		response = append(response, FlightsGetResponse{
			ID:             val.ID,
			DepCity:        val.DepCity,
			ArrCity:        val.ArrCity,
			DepTime:        val.DepTime,
			ArrTime:        val.ArrTime,
			Airplane:       val.Airplane,
			Airline:        val.Airline,
			Price:          val.Price,
			CxlSit:         val.CxlSit,
			RemainingSeats: val.LeftSeat,
		})
	}

	return c.JSONPretty(http.StatusOK, response, " ")
}
