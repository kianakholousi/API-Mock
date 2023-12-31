package handler

import (
	"errors"
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

type GetCitiesTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	e        *echo.Echo
	cities   City
	timeMock time.Time
}

func (suite *GetCitiesTestSuite) SetupSuite() {
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
	suite.cities = City{
		DB: db,
	}
	suite.timeMock = time.Date(2020, time.January, 1, 2, 3, 0, 0, time.UTC)
}

func (suite *GetCitiesTestSuite) CallHandler() (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, "/cities", strings.NewReader(""))
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := suite.cities.Get(c)

	return res, err
}

func (suite *GetCitiesTestSuite) TestGetCities_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsg := `[{"id":1,"name":"Dallas"},{"id":2,"name":"Tokyo"}]`

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "Dallas", suite.timeMock, suite.timeMock).
		AddRow(2, "Tokyo", suite.timeMock, suite.timeMock)

	suite.sqlMock.ExpectQuery("^SELECT \\* FROM `cities`$").
		RowsWillBeClosed().
		WillReturnRows(rows)

	res, err := suite.CallHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *GetCitiesTestSuite) TestGetCities_Database_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedMsg := "\"Internal Server Error\"\n"

	suite.sqlMock.ExpectQuery("^SELECT \\* FROM `cities`$").
		WillReturnError(errors.New("error"))

	res, err := suite.CallHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func TestGetCities(t *testing.T) {
	suite.Run(t, new(GetCitiesTestSuite))
}
