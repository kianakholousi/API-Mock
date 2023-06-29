package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type GetAirplanesTestSuite struct {
	suite.Suite
	sqlMock   sqlmock.Sqlmock
	e         *echo.Echo
	airplanes Airplanes
	timeMock  time.Time
}

func (suite *GetAirplanesTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	cfg := mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}

	db, err := gorm.Open(mysql.New(cfg))
	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
	suite.e = echo.New()
	suite.airplanes = Airplanes{
		DB: db,
	}
	suite.timeMock = time.Date(2020, time.January, 1, 2, 3, 0, 0, time.UTC)
}

func (suite *GetAirplanesTestSuite) CallHandler(endpoint string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, endpoint, strings.NewReader(""))
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := suite.airplanes.Get(c)

	return res, err
}

func TestGetAirplanes(t *testing.T) {
	suite.Run(t, new(GetCitiesTestSuite))
}
