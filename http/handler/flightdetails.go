package handler

import (
	"flight-data-api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type FlightDetail struct {
	DB *gorm.DB
}

type GetFlightDetailRequest struct {
	ID int32 `param:"id"`
}

type GetFlightDetailResponse struct {
	ID               int32             `json:"id"`
	DepCity          GetCitiesResponse `json:"dep_city"`
	ArrCity          GetCitiesResponse `json:"arr_city"`
	DepTime          time.Time         `json:"dep_time"`
	ArrTime          time.Time         `json:"arr_time"`
	Airline          string            `json:"airline"`
	Price            int32             `json:"price"`
	CxlSitID         int32             `json:"cxl_sit_id"`
	RemainingSeats   int32             `json:"remaining_seats"`
	FlightClass      string            `json:"flight_class"`
	BaggageAllowance int32             `json:"baggage_allowance"`
	MealService      string            `json:"meal_service"`
	GateNumber       int32             `json:"gate_number"`
}

func (f *FlightDetail) Get(ctx echo.Context) error {
	var id GetFlightDetailRequest
	if err := ctx.Bind(&id.ID); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	}

	var flight models.Flight
	err := f.DB.Model(&models.Flight{}).Select("id = ?", id.ID).Find(&flight).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, "Flight Not Found")
		}
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	responses := GetFlightDetailResponse{
		ID:               flight.ID,
		DepTime:          flight.DepTime,
		ArrTime:          flight.ArrTime,
		DepCity:          GetCitiesResponse{ID: flight.DepCity.ID, Name: flight.DepCity.Name},
		ArrCity:          GetCitiesResponse{ID: flight.DepCity.ID, Name: flight.DepCity.Name},
		Airline:          flight.Airline,
		Price:            flight.Price,
		CxlSitID:         flight.CxlSitID,
		RemainingSeats:   flight.RemainingSeats,
		FlightClass:      flight.FlightClass,
		BaggageAllowance: flight.BaggageAllowance,
		MealService:      flight.MealService,
		GateNumber:       flight.GateNumber,
	}

	return ctx.JSONPretty(http.StatusOK, responses, "")
}
