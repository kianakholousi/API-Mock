package handler

import (
	"flight-data-api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Airplanes struct {
	DB *gorm.DB
}

type GetAirplanesResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (a *Airplanes) Get(ctx echo.Context) error {
	var airplanes []models.Airplane
	err := a.DB.Debug().Find(&airplanes).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	response := make([]GetAirplanesResponse, 0, len(airplanes))
	for _, val := range airplanes {
		response = append(response, GetAirplanesResponse{ID: val.ID, Name: val.Name})
	}

	return ctx.JSON(http.StatusOK, response)
}
