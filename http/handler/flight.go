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

func (f *Flight) GetFlights(ctx echo.Context) error {
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

type GetFlightsDatesResponse struct {
	Dates []string
}

func (f *Flight) GetDates(ctx echo.Context) error {
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

	return ctx.JSONPretty(http.StatusOK, response, " ")
}

type GetFlightDetailRequest struct {
	FlightId int `json:"flight_id" validate:"required"`
}

type GetFlightDetailResponse struct {
	ID               int32                `json:"id"`
	DepCity          GetCitiesResponse    `json:"dep_city"`
	ArrCity          GetCitiesResponse    `json:"arr_city"`
	DepTime          time.Time            `json:"dep_time"`
	ArrTime          time.Time            `json:"arr_time"`
	Airplane         GetAirplanesResponse `json:"airplane"`
	Airline          string               `json:"airline"`
	Price            int32                `json:"price"`
	CxlSitID         int32                `json:"cxl_sit_id"`
	RemainingSeats   int32                `json:"remaining_seats"`
	FlightClass      string               `json:"flight_class"`
	BaggageAllowance string               `json:"baggage_allowance"`
	MealService      string               `json:"meal_service"`
	Gate             string               `json:"gate_number"`
}

func (f *Flight) GetFlightDetail(ctx echo.Context) error {
	var req GetFlightDetailRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	}

	if err := f.Validator.Struct(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	}

	var flight models.Flight
	err := f.DB.Debug().
		Model(&models.Flight{}).
		Where("id = ?", req.FlightId).
		First(&flight).
		Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return ctx.JSON(http.StatusBadRequest, "Bad Request")
	} else if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	response := GetFlightDetailResponse{
		ID:               flight.ID,
		DepCity:          GetCitiesResponse{ID: flight.DepCity.ID, Name: flight.DepCity.Name},
		ArrCity:          GetCitiesResponse{ID: flight.ArrCity.ID, Name: flight.ArrCity.Name},
		DepTime:          flight.DepTime,
		ArrTime:          flight.ArrTime,
		Airplane:         GetAirplanesResponse{ID: flight.Airplane.ID, Name: flight.Airplane.Name},
		Airline:          flight.Airline,
		Price:            flight.Price,
		CxlSitID:         flight.CxlSitID,
		RemainingSeats:   flight.RemainingSeats,
		FlightClass:      flight.FlightClass,
		BaggageAllowance: flight.BaggageAllowance,
		MealService:      flight.MealService,
		Gate:             flight.Gate,
	}

	return ctx.JSON(http.StatusOK, response)
}
