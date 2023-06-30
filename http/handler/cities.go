package handler

import (
	"flight-data-api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Cities struct {
	DB *gorm.DB
}

type GetCitiesResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (c *Cities) Get(ctx echo.Context) error {
	var cities []models.City
	err := c.DB.Debug().Find(&cities).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	response := make([]GetCitiesResponse, 0, len(cities))
	for _, val := range cities {
		response = append(response, GetCitiesResponse{ID: val.ID, Name: val.Name})
	}

	return ctx.JSON(http.StatusOK, response)
}
