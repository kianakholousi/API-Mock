package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type GetFlightsTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	e       *echo.Echo
	flight  Flight
}

func (suite *GetFlightsTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))
	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
	suite.e = echo.New()
	suite.flight = Flight{
		DB:        db,
		Validator: validator.New(),
	}
}

func (suite *GetFlightsTestSuite) CallHandler(endpoint string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, endpoint, strings.NewReader(""))
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := suite.flight.Get(c)

	return res, err
}

func TestGetFlights(t *testing.T) {
	suite.Run(t, new(GetFlightsTestSuite))
}
